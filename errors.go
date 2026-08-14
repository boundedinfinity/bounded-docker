package main

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var _ tea.Model = errorsModel{}

func newErrorModel(errs []error) tea.Model {
	columns := []table.Column{
		{Title: "#", Width: 10},
		{Title: "Error Text", Width: 120},
	}

	styles := table.DefaultStyles()

	// https://github.com/charmbracelet/lipgloss/tree/v2.0.6#rendering-tables
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(errors2Rows(errs)),
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

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Errors"
	l.SetWidth(50)
	l.ShowFilter()
	return errorsModel{
		table:  t,
		styles: lipgloss.NewStyle().Margin(1, 2),
		errs:   errs,
	}
}

type errorsModel struct {
	styles lipgloss.Style
	table  table.Model
	errs   []error
}

func errors2Rows(errs []error) []table.Row {
	rows := make([]table.Row, len(errs))
	for i := range errs {
		rows[i] = error2Row(i, errs[i])
	}
	return rows
}

func error2Row(i int, err error) table.Row {
	return table.Row{
		fmt.Sprintf("%d", i),
		err.Error(),
	}
}

func (this *errorsModel) SetErrors(errs []error) {
	this.errs = errs
}

func (this *errorsModel) ClearErrors() {
	this.errs = []error{}
}

func (this *errorsModel) Focus() {
	this.table.Focus()
}

func (this *errorsModel) Blur() {
	this.table.Blur()
}

func (this errorsModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(this.errs))
	for i := range this.errs {
		cmds[i] = func() tea.Msg { return this.errs[i] }
	}
	return tea.Batch(cmds...)
}

func (this errorsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.table.SetWidth(max(1, msg.Width))
		this.table.SetHeight(max(1, msg.Height))
	case error:
		this.errs = append(this.errs, msg)
		rows := errors2Rows(this.errs)
		this.table.SetRows(rows)
	}

	var cmd tea.Cmd
	this.table, cmd = this.table.Update(msg)
	return this, cmd
}

func (this errorsModel) View() tea.View {
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
	v := tea.NewView(baseStyle.Render(this.table.View()))
	return v
}

func createFakeErrors() []error {
	return []error{
		errors.New("Error 1"),
		errors.New("Error 2"),
		errors.New("Error 3"),
	}
}
