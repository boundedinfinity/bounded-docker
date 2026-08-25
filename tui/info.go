package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

func newInfo[T any, O any](
	tui *tui, title string, headers []string,
	options O,
	rowFn func(i int, item T) []string,
	idFn func(item T) string,
) Info {
	info := &info[T, O]{
		Title:   title,
		Items:   []T{},
		options: options,
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

type info[T any, O any] struct {
	Title       string
	Items       []T
	item        T
	options     O
	tui         *tui
	headers     []string
	header      *tview.TextView
	data        *tview.Table
	idFunc      func(T) string
	rowFunc     func(int, T) []string
	initFunc    func()
	colWidth    int
	zero        T
	selectedRow int
	selectedId  string
}

func (this *info[T, O]) Id() (string, bool) {
	row, _ := this.data.GetSelection()

	if row == this.selectedRow {
		return this.selectedId, this.selectedId != ""
	}

	idx := row - 1

	if this.idFunc != nil && idx >= 0 && idx < len(this.Items) {
		this.selectedRow = row
		item := this.Items[idx]
		this.selectedId = this.idFunc(item)
		return this.selectedId, this.selectedId != ""
	}

	return "", false
}

func (this *info[T, O]) Header() *tview.TextView {
	return this.header
}

func (this *info[T, O]) Data() *tview.Table {
	return this.data
}

func (this *info[T, O]) Focus() {
	this.tui.app.SetFocus(this.data)
}

func (this *info[T, O]) rows() ([][]string, int) {
	data := [][]string{this.headers}

	for i, item := range this.Items {
		data = append(data, this.rowFunc(i, item))
	}

	return data, len(this.Items)
}

func (this *info[T, O]) Queue(items any) {
	this.tui.wg.Go(func() {
		this.tui.app.QueueUpdateDraw(func() {
			this.Update(items)
		})
	})
}

func (this *info[T, O]) Update(items any) {
	switch v := any(items).(type) {
	case T:
		this.Set([]T{v})
	case []T:
		this.Set(v)
	case nil:
		this.Set([]T{})
	}
}

func (this *info[T, O]) Set(items []T) {
	this.Items = items
	data, count := this.rows()

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
