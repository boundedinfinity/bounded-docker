package table

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/boundedinfinity/docker-tui/tui/message"
)

var _ tea.Model = tableModel[int]{}

type tableModel[D any] struct {
	windowTitle  string
	table        table.Model
	data         []D
	data2RowFunc func(int, D) table.Row
	baseStyle    lipgloss.Style
}

func (this tableModel[D]) Init() tea.Cmd {
	return tea.Cmd(func() tea.Msg { return this.data })
}

func (this tableModel[D]) datas2Rows(summaries []D) []table.Row {
	rows := make([]table.Row, len(summaries))
	for i := range summaries {
		rows[i] = this.data2RowFunc(i, summaries[i])
	}
	return rows
}

func (this tableModel[D]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.table.SetHeight(max(1, msg.Height))
		this.table.SetWidth(max(1, msg.Width))
	case tea.FocusMsg:
		this.table.Focus()
	case tea.BlurMsg:
		this.table.Blur()
	case tea.BackgroundColorMsg:
		_ = msg.IsDark()
	case message.ClearMsg:
		this.data = []D{}
	case D:
		this.data = append(this.data, msg)
	case []D:
		this.data = msg
	}

	this.table, cmd = this.table.Update(msg)
	return this, cmd
}

func (this tableModel[D]) View() tea.View {
	rows := this.datas2Rows(this.data)
	this.table.SetRows(rows)
	v := tea.NewView(this.baseStyle.Render(this.table.View()))
	v.WindowTitle = this.windowTitle

	return v
}
