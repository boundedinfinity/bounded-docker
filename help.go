package main

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/boundedinfinity/docker-tui/state"
)

func newHelp(state *state.Machine) helpModel {
	return helpModel{
		help:        help.New(),
		inputStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#1f2d4b")),
		widgetStyle: lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
		state:       state,
		current:     state.Current,
	}
}

type helpModel struct {
	help        help.Model
	inputStyle  lipgloss.Style
	widgetStyle lipgloss.Style
	quitting    bool
	state       *state.Machine
	current     *state.State
}

func (this helpModel) Init() tea.Cmd {
	return nil
}

func (this helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.help.SetWidth(msg.Width)
	case tea.KeyPressMsg:
		if state, ok := this.state.Next(msg); ok {
			this.current = state
		}
	}

	return this, nil
}

func (this helpModel) View() tea.View {
	b := this.current.HelpView()
	v := this.help.ShortHelpView(b)
	v = this.widgetStyle.Render(v)
	return tea.NewView(v)
}
