package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

func newInfo[T any](tui *tui, title string, headers []string, fn func(i int, item T) []string) Info {
	info := &info[T]{
		Title:        title,
		Items:        []T{},
		headers:      headers,
		tui:          tui,
		item2RowFunc: fn,
		selectedRow:  0,
	}

	info.header = tview.NewTextView()

	info.data = tview.NewTable().
		SetBorders(true).
		SetFixed(1, len(headers)).
		SetEvaluateAllRows(true).
		SetSelectable(true, false)

	return info
}

type Info interface {
	Update(items any)
	Focus()
	Header() *tview.TextView
	Data() *tview.Table
	Init()
}

type info[T any] struct {
	Title        string
	Items        []T
	tui          *tui
	headers      []string
	header       *tview.TextView
	data         *tview.Table
	item2RowFunc func(int, T) []string
	colWidth     int
	selectedRow  int
}

func (this *info[T]) Init() {
	this.data.SetSelectedFunc(func(row, _ int) { this.selectedRow = row })
	this.update([]T{})
}

func (this *info[T]) Update(items any) {
	switch v := any(items).(type) {
	case T:
		this.Append(v)
	case []T:
		this.Set(v)
	case nil:
		this.Set([]T{})
	}
}

func (this *info[T]) Header() *tview.TextView {
	return this.header
}

func (this *info[T]) Data() *tview.Table {
	return this.data
}

func (this *info[T]) Append(item T) {
	this.Set(append(this.Items, item))
}

func (this *info[T]) Set(items []T) {
	this.tui.app.QueueUpdateDraw(func() {
		this.update(items)
	})
}

func (this *info[T]) Focus() {
	this.tui.app.SetFocus(this.data)
}

func (this *info[T]) update(items []T) {
	this.Items = items
	count := len(this.Items)

	data := [][]string{this.headers}
	for i, item := range this.Items {
		data = append(data, this.item2RowFunc(i, item))
	}

	this.header.Clear()
	fmt.Fprintf(this.header, "%s [%d]", this.Title, count)

	this.data.Clear()
	for row, vals := range data {
		for col, text := range vals {
			cell := tview.NewTableCell(text)
			cell.SetExpansion(1)

			if row == 0 {
				cell.SetSelectable(false)
			}

			this.data.SetCell(row, col, cell)
		}
	}
}
