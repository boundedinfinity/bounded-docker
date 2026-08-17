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
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var _ tea.Model = appModel{}

func newApp() (appModel, error) {
	cmodel, cmenu := newContainersModel(createFakeContainers())
	imodel, imenu := newImageModel(createFakeImages())
	emodel, emenu := newErrorModel(createFakeErrors())
	views := map[string]tea.Model{
		"root":       newWelcome(),
		"containers": cmodel,
		"images":     imodel,
		"errors":     emodel,
	}

	machine, err := state.New(state.DefaultConfig(), views)

	if err != nil {
		return appModel{}, err
	}

	m := appModel{
		state: machine,
		menu:  menu.New(machine, cmenu, imenu, emenu),
		help:  newHelp(machine),
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),
	}

	m.view = m.state.Start()

	return m, nil
}

type appModel struct {
	style lipgloss.Style
	// cancel context.CancelFunc
	menu  tea.Model
	help  tea.Model
	state *state.Machine
	view  tea.Model
}

func (this appModel) Init() tea.Cmd {
	return nil
}

func (this appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	broadcast := func(msg tea.Msg) {
		this.help, cmd = this.help.Update(msg)
		cmds = append(cmds, cmd)
		this.menu, cmd = this.menu.Update(msg)
		cmds = append(cmds, cmd)
		cmd = this.state.Broadcast(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newSize := this.size(msg)
		broadcast(newSize)
	case []container.Summary, []image.Summary, error:
		broadcast(msg)
	case tea.KeyPressMsg:
		model, cmd, ok := this.state.Update(msg)
		if ok {
			return model, cmd
		}
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
	join := lipgloss.JoinVertical(
		lipgloss.Top,
		this.menu.View().Content,
		this.view.View().Content,
		this.help.View().Content,
	)

	v := tea.NewView(this.style.Render(join))

	v.AltScreen = true
	v.WindowTitle = "Bounded Docker"
	return v
}

func main() {
	var err error
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	d, err := docker.NewDocker(wg, ctx)
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		os.Exit(1)
	}

	m, err := newApp()
	if err != nil {
		fmt.Println("Error creating app:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithContext(ctx))

	go func() {
		<-sigCh
		cancel()
	}()

	wg.Go(func() {
		if _, err = p.Run(); err != nil {
			cancel()
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

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
