package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type welcomeModel struct {
	textStyle lipgloss.Style
	boxStyle  lipgloss.Style
}

func newWelcome() welcomeModel {
	return welcomeModel{
		textStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#2a1981")),
		boxStyle:  lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
	}
}

func (this welcomeModel) Init() tea.Cmd {
	return nil
}

func (this welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return this, nil
}

func (this welcomeModel) View() tea.View {
	v := "Welcome to Bounded Docker TUI!\n\n"
	v = this.textStyle.Render(v)
	v = this.boxStyle.Render(v)

	return tea.NewView(v)
}
