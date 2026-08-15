package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/boundedinfinity/docker-tui/state"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
)

var _ tea.Model = menuModel{}

func newMenu(state *state.Machine) tea.Model {
	return menuModel{
		summaries: 0,
		images:    0,
		errs:      0,
		state:     state,
	}
}

type menuModel struct {
	summaries int
	images    int
	errs      int
	state     *state.Machine
}

func (this menuModel) Init() tea.Cmd {
	return nil
}

func (this menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		_ = msg.Width
		_ = msg.Height
	case []error:
		this.errs = len(msg)
	case []container.Summary:
		this.summaries = len(msg)
	case []image.Summary:
		this.images = len(msg)
	}

	return this, nil
}

func (this menuModel) View() tea.View {
	styles := lipgloss.NewStyle().Padding(0, 10)

	style := func(format string, a ...any) string {
		text := fmt.Sprintf(format, a...)
		return styles.Render(text)
	}

	join := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Padding(0, 0, 0, 1).Render("Bounded Docker:"),
		style("Containers [%d]", this.summaries),
		style("Images [%d]", this.images),
		style("Errors [%d]", this.errs),
	)

	return tea.NewView(join)
}
