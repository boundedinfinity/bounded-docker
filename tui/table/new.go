package table

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func New[D any](title string, columns []table.Column, data2RowFunc func(int, D) table.Row, data []D) tea.Model {
	m := tableModel[D]{
		windowTitle:  title,
		data:         data,
		data2RowFunc: data2RowFunc,
		baseStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),
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
