package main

import (
	"github.com/boundedinfinity/docker-tui/state"
	"github.com/boundedinfinity/docker-tui/tui"
)

// var _ tea.Model = appModel{}

// func newApp() (appModel, error) {
// 	cmodel, cmenu := newContainersModel(createFakeContainers())
// 	imodel, imenu := newImageModel(createFakeImages())
// 	emodel, emenu := newErrorModel(createFakeErrors())
// 	views := map[string]tea.Model{
// 		"root":       newWelcome(),
// 		"containers": cmodel,
// 		"images":     imodel,
// 		"errors":     emodel,
// 	}

// 	m := appModel{
// 		state: machine,
// 		menu:  menu.New(machine, cmenu, imenu, emenu),
// 		help:  newHelp(machine),
// 		views: views,
// 		style: lipgloss.NewStyle().
// 			BorderStyle(lipgloss.NormalBorder()).
// 			BorderForeground(lipgloss.Color("240")),
// 	}

// 	m.view = m.state.Start()

// 	return m, nil
// }

// type appModel struct {
// 	style lipgloss.Style
// 	menu  tea.Model
// 	help  tea.Model
// 	state *state.Machine
// 	view  tea.Model
// 	views map[string]tea.Model
// }

// func (this appModel) Init() tea.Cmd {
// 	return nil
// }

// func (this appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmds []tea.Cmd
// 	var cmd tea.Cmd

// 	update := func(view tea.Model, msg tea.Msg) tea.Model {
// 		view, cmd = view.Update(msg)
// 		cmds = append(cmds, cmd)
// 		return view
// 	}

// 	this.view = update(this.state.Update(msg))

// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		newSize := this.size(msg)
// 		this.menu = update(this.menu, newSize)
// 		this.help = update(this.help, newSize)
// 		this.view = update(this.view, newSize)
// 	case error:
// 		this.menu = update(this.menu, msg)
// 		this.views["errors"] = update(this.views["errors"], msg)
// 	}

// 	return this, tea.Batch(cmds...)
// }

// func (this appModel) size(msg tea.WindowSizeMsg) tea.WindowSizeMsg {
// 	borderW, borderH := lipgloss.Size(this.style.Render(""))
// 	_, menuH := lipgloss.Size(this.menu.View().Content)
// 	_, helpH := lipgloss.Size(this.help.View().Content)
// 	return tea.WindowSizeMsg{
// 		Width:  msg.Width - borderW*2,
// 		Height: msg.Height - borderH*4 - menuH - helpH,
// 	}
// }

// func (this appModel) View() tea.View {
// 	join := lipgloss.JoinVertical(
// 		lipgloss.Top,
// 		this.menu.View().Content,
// 		this.view.View().Content,
// 		this.help.View().Content,
// 	)

// 	v := tea.NewView(this.style.Render(join))

// 	v.AltScreen = true
// 	v.WindowTitle = "Bounded Docker"
// 	return v
// }

func main() {
	// var err error
	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, os.Interrupt)

	// ctx, cancel := context.WithCancel(context.Background())
	// wg := &sync.WaitGroup{}

	// d, err := docker.NewDocker(wg, ctx)
	// if err != nil {
	// 	fmt.Println("Error creating Docker client:", err)
	// 	os.Exit(1)
	// }

	// wg.Go(func() {
	// 	if _, err = p.Run(); err != nil {
	// 		cancel()
	// 	}
	// })

	// wg.Go(func() {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		case err := <-d.ErrCh:
	// 			p.Send(err)
	// 		case summary := <-d.ContainersCh:
	// 			p.Send(summary)
	// 		case images := <-d.ImagesCh:
	// 			p.Send(images)
	// 		}
	// 	}
	// })

	// d.Init()
	// wg.Wait()

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	os.Exit(1)
	// }

	// go func() {
	// 	<-sigCh
	// 	cancel()
	// }()

	machine, err := state.New(state.DefaultConfig())

	if err != nil {
		panic(err)
	}

	tui := tui.NewTui(machine)

	if err := tui.Run(); err != nil {
		panic(err)
	}
}
