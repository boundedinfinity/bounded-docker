package main

import (
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/moby/moby/api/types/container"
)

var _ tea.Model = summaryModel{}

func newSummaryModel() tea.Model {
	columns := []table.Column{
		{Title: "ID"},
		{Title: "Name"},
		{Title: "Image"},
		{Title: "Command"},
		{Title: "Status"},
	}

	// rows := []table.Row{}
	rows := []table.Row{
		{"1234567890", "my-container", "my-image", "/bin/bash", "Up 5 minutes"},
	}

	// https://github.com/charmbracelet/lipgloss/tree/v2.0.6#rendering-tables
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(0),
		table.WithWidth(0),
	)

	t.Help.ShowAll = false

	return summaryModel{
		// styles: getSummaryStyles(),
		table: t,
	}
}

func getSummaryStyles() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
}

type summaryModel struct {
	// styles    lipgloss.Style
	table     table.Model
	summaries []container.Summary
}

func (this *summaryModel) SetSummaries(summaries []container.Summary) {
	this.summaries = summaries
}

func (this *summaryModel) ClearSummaries() {
	this.summaries = []container.Summary{}
}

func (this summaryModel) Init() tea.Cmd {
	return nil
}

func summary2Row(summary container.Summary) table.Row {
	names := strings.Join(summary.Names, ", ")
	return table.Row{
		summary.ID,
		names,
		summary.Image,
		summary.Command,
		summary.Status,
	}
}

func (this summaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.table.SetHeight(msg.Height - 10)
		this.table.SetWidth(msg.Width - 50)
	case tea.BackgroundColorMsg:
		_ = msg.IsDark()
	case []container.Summary:
		this.summaries = msg
		rows := make([]table.Row, len(this.summaries))
		for i := range this.summaries {
			rows[i] = summary2Row(this.summaries[i])
		}
		this.table.SetRows(rows)
	}

	this.table, cmd = this.table.Update(msg)
	return this, cmd
}

func (this summaryModel) View() tea.View {
	v := tea.NewView(this.table.View())
	// v.WindowTitle = "Containers"

	return v
}

// func createContainerTable(pageWidth int, columns []table.Column, rows []table.Row) table.Model {
// 	pagePadding := 20
// 	columnPadding := 5

// 	tableWidth := max(0, pageWidth-pagePadding) // account for borders
// 	columnWidth := int(tableWidth/len(columns)) + columnPadding
// 	widths := make([]int, len(columns))

// 	for i := range columns {
// 		widths[i] = max(widths[i], columns[i].Width)
// 	}

// 	for r := range rows {
// 		for i := range rows[r] {
// 			widths[i] = max(widths[i], len(rows[r][i]))
// 		}
// 	}

// 	for i := range columns {
// 		columns[i].Width = min(columnWidth, widths[i]+columnPadding)
// 	}

// 	t := table.New(
// 		table.WithColumns(columns),
// 		table.WithRows(rows),
// 		table.WithFocused(true),
// 		table.WithHeight(7),
// 		table.WithWidth(tableWidth),
// 	)

// 	s := table.DefaultStyles()

// 	s.Header = s.Header.
// 		BorderStyle(lipgloss.NormalBorder()).
// 		BorderForeground(lipgloss.Color("240")).
// 		BorderBottom(true).
// 		Bold(false)

// 	s.Selected = s.Selected.
// 		Foreground(lipgloss.Color("229")).
// 		Background(lipgloss.Color("57")).
// 		Bold(false)

// 	t.SetStyles(s)

// 	return t
// }
