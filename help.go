package main

import (
	"strings"

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
	var c []string

	if th := this.state.Current.TransisionHelp(); len(th) > 0 {
		tv := this.help.ShortHelpView(th)
		c = append(c, tv)
	}

	if nh := this.state.Current.NavigationHelp(); len(nh) > 0 {
		nv := this.help.ShortHelpView(nh)
		c = append(c, nv)
	}

	v := lipgloss.JoinVertical(lipgloss.Top, strings.Join(c, "\n\n"))

	// b := this.state.Current.HelpView2()
	// v := this.help.FullHelpView(b)

	v = this.borderStyle.Render(v)
	return tea.NewView(v)
}
