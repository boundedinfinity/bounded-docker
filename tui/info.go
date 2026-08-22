package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

func newInfo[T any](
	tui *tui, title string, headers []string,
	rowFn func(i int, item T) []string,
	idFn func(item T) string,
) Info {
	info := &info[T]{
		Title:   title,
		Items:   []T{},
		headers: headers,
		tui:     tui,
		rowFunc: rowFn,
		idFunc:  idFn,
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
	Queue(items any)
	Update(items any)
	Focus()
	Header() *tview.TextView
	Data() *tview.Table
	Id() (string, bool)
}

type info[T any] struct {
	Title    string
	Items    []T
	item     T
	tui      *tui
	headers  []string
	header   *tview.TextView
	data     *tview.Table
	idFunc   func(T) string
	rowFunc  func(int, T) []string
	initFunc func()
	colWidth int
	zero     T
}

func (this *info[T]) Id() (string, bool) {
	row, _ := this.data.GetSelection()
	idx := row - 1

	if this.idFunc != nil && idx >= 0 && idx < len(this.Items) {
		item := this.Items[idx]
		id := this.idFunc(item)
		return id, id != ""
	}

	return "", false
}

func (this *info[T]) Header() *tview.TextView {
	return this.header
}

func (this *info[T]) Data() *tview.Table {
	return this.data
}

func (this *info[T]) Focus() {
	this.tui.app.SetFocus(this.data)
}

func (this *info[T]) rows() ([][]string, int) {
	data := [][]string{this.headers}

	for i, item := range this.Items {
		data = append(data, this.rowFunc(i, item))
	}

	return data, len(this.Items)
}

func (this *info[T]) Queue(items any) {
	this.tui.wg.Go(func() {
		this.tui.app.QueueUpdate(func() {
			this.Update(items)
		})
	})
}

func (this *info[T]) Update(items any) {
	switch v := any(items).(type) {
	case T:
		this.Set([]T{v})
	case []T:
		this.Set(v)
	case nil:
		this.Set([]T{})
	}
}

func (this *info[T]) Set(items []T) {
	this.Items = items
	data, count := this.rows()

	// this.header.Clear()
	fmt.Fprintf(this.header, "%s [%d]", this.Title, count)

	// this.data.Clear()
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
