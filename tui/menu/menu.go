package menu

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/boundedinfinity/docker-tui/state"
)

var _ tea.Model = model{}

func New(state *state.Machine, items ...tea.Model) tea.Model {
	return model{
		baseStyle:  lipgloss.NewStyle().Background(lipgloss.Color("#fff")),
		titleStyle: lipgloss.NewStyle().Padding(0, 0, 0, 1),
		itemStyle:  lipgloss.NewStyle().Padding(0, 10),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#211f1f")).
			Background(lipgloss.Color("#6f6c6c")).
			Padding(0, 10),
		items: items,
		state: state,
	}
}

type model struct {
	baseStyle     lipgloss.Style
	titleStyle    lipgloss.Style
	itemStyle     lipgloss.Style
	selectedStyle lipgloss.Style
	items         []tea.Model
	state         *state.Machine
}

func (this model) Init() tea.Cmd {
	return nil
}

func (this model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	var cmd tea.Cmd

	for i := range this.items {
		this.items[i], cmd = this.items[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return this, tea.Batch(cmds...)
}

func (this model) View() tea.View {
	items := []string{this.titleStyle.Render("Bounded Docker:")}

	for _, item := range this.items {
		if mitem, ok := item.(MenuItem); ok {
			v := fmt.Sprintf("%s [%d]", mitem.Title(), mitem.Count())
			s := this.itemStyle

			if this.state.Current.Name == mitem.Title() {
				s = this.selectedStyle
			}

			v = s.Render(v)
			items = append(items, v)
		}
	}

	join := lipgloss.JoinHorizontal(
		lipgloss.Left,
		items...,
	)

	return tea.NewView(join)
}
