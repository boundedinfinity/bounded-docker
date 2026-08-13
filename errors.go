package main

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var _ tea.Model = errorsModel{}

func newErrorModel() tea.Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Errors"
	l.ShowFilter()
	return errorsModel{
		list:   l,
		errs:   make([]error, 0),
		styles: lipgloss.NewStyle().Margin(1, 2),
	}
}

type errItem struct {
	err error
}

func (this errItem) Title() string {
	return this.err.Error()
}

func (this errItem) FilterValue() string {
	return this.err.Error()
}

type errorsModel struct {
	styles lipgloss.Style
	list   list.Model
	errs   []error
}

func (this *errorsModel) SetErrors(errs []error) {
	this.errs = errs
}

func (this *errorsModel) ClearErrors() {
	this.errs = []error{}
}

func (this errorsModel) Init() tea.Cmd {
	return nil
}

func (this errorsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		this.errs = append(this.errs, msg)
	}

	var cmd tea.Cmd
	this.list, cmd = this.list.Update(msg)
	return this, cmd
}

func (this errorsModel) View() tea.View {
	items := make([]list.Item, len(this.errs))
	for i, err := range this.errs {
		items[i] = errItem{err: err}
	}

	this.list.SetItems(items)

	v := tea.NewView(this.styles.Render(this.list.View()))
	return v
}
