package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/boundedinfinity/docker-tui/docker"
	"github.com/boundedinfinity/docker-tui/state"
	"github.com/boundedinfinity/docker-tui/tui/menu"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var _ tea.Model = appModel{}

func newApp(ctx context.Context, cancel context.CancelFunc) (appModel, error) {
	machine, err := state.New(state.DefaultConfig())

	if err != nil {
		return appModel{}, err
	}

	ct, cm := newContainersModel(createFakeContainers())
	it, im := newImageModel(createFakeImages())
	et, em := newErrorModel(createFakeErrors())

	m := appModel{
		current: "root",
		ctx:     ctx,
		cancel:  cancel,
		state:   machine,
		menu:    menu.New(machine, cm, im, em),
		help:    newHelp(machine),
		pages: map[string]tea.Model{
			"root":       newWelcome(),
			"containers": ct,
			"images":     it,
			"errors":     et,
		},
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),
	}

	return m, nil
}

type appModel struct {
	style   lipgloss.Style
	ctx     context.Context
	cancel  context.CancelFunc
	menu    tea.Model
	help    tea.Model
	pages   map[string]tea.Model
	current string
	state   *state.Machine
}

func (this appModel) Init() tea.Cmd {
	return nil
}

func (this appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	helper := func(m tea.Model, msg2 tea.Msg) tea.Model {
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
			this.pages[i] = helper(this.pages[i], newSize)
		}
		return this, tea.Batch(cmds...)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			this.cancel()
			return this, nil
		default:
			if next, ok := this.state.Next(msg); ok {
				if m, ok := this.pages[next.Id]; ok {
					m := helper(m, tea.FocusMsg{})
					this.pages[next.Id] = m
					this.current = next.Id
				}
			}
		}
	}

	this.help = helper(this.help, msg)
	this.menu = helper(this.menu, msg)

	for i := range this.pages {
		this.pages[i] = helper(this.pages[i], msg)
	}

	return this, tea.Batch(cmds...)
}

func (this appModel) size(msg tea.WindowSizeMsg) tea.WindowSizeMsg {
	borderW, borderH := lipgloss.Size(this.style.Render(""))
	_, menuH := lipgloss.Size(this.menu.View().Content)
	_, helpH := lipgloss.Size(this.help.View().Content)
	return tea.WindowSizeMsg{
		Width:  msg.Width - borderW*2,
		Height: msg.Height - borderH*4 - menuH - helpH,
	}
}

func (this appModel) View() tea.View {
	current, _ := this.pages[this.current]

	join := lipgloss.JoinVertical(
		lipgloss.Top,
		this.menu.View().Content,
		current.View().Content,
		this.help.View().Content,
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
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		os.Exit(1)
	}

	signal.Notify(sigCh, os.Interrupt)

	m, err := newApp(ctx, cancel)
	if err != nil {
		fmt.Println("Error creating app:", err)
		os.Exit(1)
	}

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
			case summary := <-d.ContainersCh:
				p.Send(summary)
			case images := <-d.ImagesCh:
				p.Send(images)
			}
		}
	})

	d.Init()
	wg.Wait()
}
