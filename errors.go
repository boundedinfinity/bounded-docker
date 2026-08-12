package main

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type errItem struct {
	err error
}

func (this errItem) Title() string {
	return this.err.Error()
}

func (this errItem) FilterValue() string {
	return this.err.Error()
}

type errModel struct {
	list list.Model
	errs []error
}

func (m errModel) Init() tea.Cmd {
	return nil
}

func (m errModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		m.errs = append(m.errs, msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m errModel) View() tea.View {
	items := make([]list.Item, len(m.errs))
	for i, err := range m.errs {
		items[i] = errItem{err: err}
	}

	m.list.SetItems(items)

	v := tea.NewView(docStyle.Render(m.list.View()))
	v.AltScreen = true
	return v
}

func createErrorView() errModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Errors"
	return errModel{
		list: l,
		errs: make([]error, 0),
	}
}
