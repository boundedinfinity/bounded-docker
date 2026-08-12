package main

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

func createContainerTable(pageWidth int, columns []table.Column, rows []table.Row) table.Model {
	pagePadding := 20
	columnPadding := 5

	tableWidth := max(0, pageWidth-pagePadding) // account for borders
	columnWidth := int(tableWidth/len(columns)) + columnPadding
	widths := make([]int, len(columns))

	for i := range columns {
		widths[i] = max(widths[i], columns[i].Width)
	}

	for r := range rows {
		for i := range rows[r] {
			widths[i] = max(widths[i], len(rows[r][i]))
		}
	}

	for i := range columns {
		columns[i].Width = min(columnWidth, widths[i]+columnPadding)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(tableWidth),
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

	return t
}
