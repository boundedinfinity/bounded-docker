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

var _ tea.Model = &appModel{}

func newApp(ctx context.Context, cancel context.CancelFunc) *appModel {
	m := &appModel{
		ctx:         ctx,
		cancel:      cancel,
		summaries:   newSummaryModel(),
		errorModels: newErrorModel(),
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
	summaries   tea.Model
	errorModels tea.Model
}

func (this *appModel) Init() tea.Cmd {
	return nil
}

func (this *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	select {
	case <-this.ctx.Done():
		return this, tea.Quit
	default:
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// this.windowWidth = msg.Width
		// this.table = createContainerTable(msg.Width, this.columns, this.rows)
	case tea.KeyPressMsg:
		switch msg.String() {
		// case "esc":
		// 	if this.table.Focused() {
		// 		this.table.Blur()
		// 	} else {
		// 		this.table.Focus()
		// 	}
		case "q", "ctrl+c":
			this.cancel()
			return this, nil
		case "enter":
			// return this, tea.Batch(
			// 	tea.Printf("Let's go to %s!", this.table.SelectedRow()[1]),
			// )
		}
	case error:
		return this, nil
	}

	this.summaries, cmd = this.summaries.Update(msg)
	return this, cmd
}

func (this appModel) View() tea.View {
	// join := lipgloss.JoinHorizontal(
	// 	0,
	// 	this.summaries.View().Content,
	// )

	v := tea.NewView(this.summaries.View().Content)

	v.AltScreen = true
	v.WindowTitle = "Bounded Docker"
	return v

	// return tea.NewView("Ths is a test")
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
