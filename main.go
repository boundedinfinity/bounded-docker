package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/boundedinfinity/docker-tui/docker"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var _ tea.Model = appModel{}

type page int

const (
	summaries page = iota
	errs
)

func newApp(ctx context.Context, cancel context.CancelFunc) appModel {
	m := appModel{
		ctx:         ctx,
		cancel:      cancel,
		menu:        newMenu(),
		currentPage: summaries,
		pages: []tea.Model{
			newSummaryModel(createFakeSummaries()),
			newErrorModel(createFakeErrors()),
		},
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),
	}

	return m
}

type appModel struct {
	style       lipgloss.Style
	ctx         context.Context
	cancel      context.CancelFunc
	menu        tea.Model
	pages       []tea.Model
	currentPage page
}

func (this appModel) Init() tea.Cmd {
	return nil
}

func (this appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	helper := func(p page, m tea.Model, msg2 tea.Msg) tea.Model {
		if f, ok := m.(Focusable); ok {
			if p == this.currentPage {
				f.Focus()
			} else {
				f.Blur()
			}
		}

		m, cmd := m.Update(msg2)
		cmds = append(cmds, cmd)
		return m
	}

	select {
	case <-this.ctx.Done():
		return this, tea.Quit
	default:
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newSize := this.size(msg)
		for i := range this.pages {
			this.pages[i] = helper(page(i), this.pages[i], newSize)
		}
		return this, tea.Batch(cmds...)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "c":
			this.currentPage = summaries
		case "e":
			this.currentPage = errs
		case "q", "ctrl+c":
			this.cancel()
			return this, nil
		}
	}

	this.menu = helper(-1, this.menu, msg)
	for i := range this.pages {
		this.pages[i] = helper(page(i), this.pages[i], msg)
	}

	return this, tea.Batch(cmds...)
}

func (this appModel) size(msg tea.WindowSizeMsg) tea.WindowSizeMsg {
	borderW, borderH := lipgloss.Size(this.style.Render(""))
	_, menuH := lipgloss.Size(this.menu.View().Content)
	return tea.WindowSizeMsg{
		Width:  msg.Width - borderW*2,
		Height: msg.Height - borderH*4 - menuH,
	}
}

func (this appModel) View() tea.View {
	m := this.pages[this.currentPage]

	join := lipgloss.JoinVertical(
		lipgloss.Top,
		this.menu.View().Content,
		m.View().Content,
	)

	v := tea.NewView(this.style.Render(join))

	v.AltScreen = true
	v.WindowTitle = "Bounded Docker"
	return v
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	sigCh := make(chan os.Signal, 1)
	d, err := docker.NewDocker(wg, ctx)

	signal.Notify(sigCh, os.Interrupt)

	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		os.Exit(1)
	}

	m := newApp(ctx, cancel)
	p := tea.NewProgram(m)

	go func() {
		<-sigCh
		cancel()
		p.Quit()
	}()

	wg.Go(func() {
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			cancel()
			os.Exit(1)
		}
	})

	wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-d.ErrCh:
				p.Send(err)
			case summary := <-d.SummaryCh:
				p.Send(summary)
			}
		}
	})

	d.GetSummary()

	wg.Wait()
}
