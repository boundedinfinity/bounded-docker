package details

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

type Item struct {
	Name  string
	Value string
}

func Details[T any](data2View func(T) []Item) *DetailsView[T] {
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Value", Width: 30},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithHeight(0),
		table.WithWidth(0),
		table.WithStyles(table.DefaultStyles()),
	)

	return &DetailsView[T]{
		data2Details: data2View,
		table:        t,
	}
}

var _ tea.Model = &DetailsView[any]{}

type DetailsView[T any] struct {
	data         T
	data2Details func(T) []Item
	table        table.Model
}

func (this *DetailsView[T]) Init() tea.Cmd {
	return nil
}

func (this *DetailsView[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case T:
		this.data = msg
	}

	return this, nil
}

func (this *DetailsView[T]) View() tea.View {
	items := this.data2Details(this.data)
	rows := make([]table.Row, 0, len(items))

	for i := range items {
		rows = append(rows, table.Row{items[i].Name, items[i].Value})
	}

	this.table.SetRows(rows)
	return tea.NewView(this.table.View())
}
