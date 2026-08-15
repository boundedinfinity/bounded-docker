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
		borderStyle: lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
		state:       state,
	}
}

type helpModel struct {
	help        help.Model
	borderStyle lipgloss.Style
	state       *state.Machine
}

func (this helpModel) Init() tea.Cmd {
	return nil
}

func (this helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.help.SetWidth(msg.Width)
	}

	return this, nil
}

func (this helpModel) View() tea.View {
	b := this.state.Current.HelpView()
	v := this.help.ShortHelpView(b)
	v = this.borderStyle.Render(v)
	return tea.NewView(v)
}
