package main

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func newTableModel[D any](title string, columns []table.Column, data2Row func(int, D) table.Row, data []D) tea.Model {
	m := tableModel[D]{
		windowTitle: title,
		data:        data,
		data2Row:    data2Row,
	}
	styles := table.DefaultStyles()

	// https://github.com/charmbracelet/lipgloss/tree/v2.0.6#rendering-tables
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(m.datas2Rows(data)),
		table.WithHeight(0),
		table.WithWidth(0),
		table.WithStyles(styles),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	t.Help.ShowAll = false
	m.windowTitle = title
	m.table = t

	return m
}

var _ tea.Model = tableModel[int]{}

type tableModel[D any] struct {
	windowTitle string
	table       table.Model
	data        []D
	data2Row    func(int, D) table.Row
}

func (this tableModel[D]) Init() tea.Cmd {
	return tea.Cmd(func() tea.Msg { return this.data })
}

func (this tableModel[D]) datas2Rows(summaries []D) []table.Row {
	rows := make([]table.Row, len(summaries))
	for i := range summaries {
		rows[i] = this.data2Row(i, summaries[i])
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
	case ClearMsg:
		this.data = []D{}
		rows := this.datas2Rows(this.data)
		this.table.SetRows(rows)
	case D:
		this.data = append(this.data, msg)
		rows := this.datas2Rows(this.data)
		this.table.SetRows(rows)
	case []D:
		this.data = msg
		rows := this.datas2Rows(msg)
		this.table.SetRows(rows)
	}

	this.table, cmd = this.table.Update(msg)
	return this, cmd
}

func (this tableModel[D]) View() tea.View {
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	v := tea.NewView(baseStyle.Render(this.table.View()))
	v.WindowTitle = this.windowTitle

	return v
}
