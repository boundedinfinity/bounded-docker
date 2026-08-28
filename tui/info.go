package tui

import (
	"fmt"

	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/rivo/tview"
)

type Info interface {
	Queue(any)
	Update(any)
	Header() *tview.TextView
	Data() *tview.Table
	Id() (string, bool)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

var (
	ErrInfoList   = fmt.Errorf("info list")
	errInfoListFn = errorer.Func(ErrInfoList)
)

func newInfoList[T any, O any](
	tui *tui, title string, options O,
	rowsFn func(items []T) [][]string,
	idFn func(item T) string,
) Info {
	info := &info[T, O]{
		Title:   title,
		Items:   []T{},
		options: options,
		tui:     tui,
		rowsFn:  rowsFn,
		idFunc:  idFn,
	}

	info.header = tview.NewTextView()

	info.data = tview.NewTable().
		SetBorders(true).
		SetFixed(1, 1).
		SetEvaluateAllRows(true).
		SetSelectable(true, false)

	info.Update([]T{})

	return info
}

type info[T any, O any] struct {
	Title       string
	Items       []T
	item        T
	options     O
	tui         *tui
	header      *tview.TextView
	data        *tview.Table
	idFunc      func(T) string
	rowsFn      func([]T) [][]string
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

func (this *info[T, O]) ScrollToId(id string) {
	for row, item := range this.Items {
		if this.idFunc(item) == id {
			this.data.Select(row+1, 0)
			break
		}
	}
}

func (this *info[T, O]) Header() *tview.TextView {
	return this.header
}

func (this *info[T, O]) Data() *tview.Table {
	return this.data
}

func (this *info[T, O]) rows() ([][]string, int) {
	return this.rowsFn(this.Items), len(this.Items)
}

func (this *info[T, O]) Queue(items any) {
	this.tui.queueDraw(func() {
		this.Update(items)
	})
}

func (this *info[T, O]) Update(a any) {
	switch items := any(a).(type) {
	case T:
		this.Items = []T{items}
	case []T:
		this.Items = items
	case nil:
		this.Items = []T{}
	default:
		panic(errInfoListFn(
			"invalid type: %T, expected %T or []%T",
			a, this.zero, this.zero,
		))
	}

	rows, count := this.rows()

	this.header.Clear()
	fmt.Fprintf(this.header, "%s [%d]", this.Title, count)

	this.data.Clear()
	this.data.SetFixed(1, len(rows[0]))
	for row, vals := range rows {
		for col, text := range vals {
			cell := tview.NewTableCell(text)
			cell.SetExpansion(1)

			if row == 0 {
				cell.SetSelectable(false)
			}

			this.data.SetCell(row, col, cell)
		}
	}

	this.data.ScrollToBeginning()
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

var (
	ErrInfoItem   = fmt.Errorf("info item")
	errInfoItemFn = errorer.Func(ErrInfoItem)
)

func newInfoItem[T any](
	tui *tui, title string,
	rowsFunc func(T) [][]string,
	idFn func(item T) string,
) Info {
	info := &infoItem[T]{
		title:    title,
		tui:      tui,
		rowsFunc: rowsFunc,
		idFunc:   idFn,
	}

	info.header = tview.NewTextView()

	info.data = tview.NewTable().
		SetBorders(true).
		SetFixed(1, 3).
		SetEvaluateAllRows(true).
		SetSelectable(true, false)

	return info
}

type infoItem[T any] struct {
	item     T
	tui      *tui
	title    string
	header   *tview.TextView
	data     *tview.Table
	idFunc   func(T) string
	rowsFunc func(T) [][]string
}

func (this *infoItem[T]) Id() (string, bool) {
	if this.idFunc == nil {
		panic(errInfoItemFn("idFunc is nil"))
	}

	return this.idFunc(this.item), true
}

func (this *infoItem[T]) Header() *tview.TextView {
	return this.header
}

func (this *infoItem[T]) Data() *tview.Table {
	return this.data
}

func (this *infoItem[T]) Queue(a any) {
	this.tui.queueDraw(func() {
		this.Update(a)
	})
}

func (this *infoItem[T]) Update(a any) {
	if this.rowsFunc == nil {
		panic(errInfoItemFn("rowsFunc is nil"))
	}

	if item, ok := a.(T); ok {
		this.item = item
	} else {
		var zero T
		panic(errInfoItemFn("invalid type: %T, expected: %T", a, zero))
	}

	id, _ := this.Id()
	title := fmt.Sprintf("%s [%s]", this.title, id)
	this.header.Clear()
	fmt.Fprintf(this.header, "%s", title)

	rows := this.rowsFunc(this.item)
	this.data.Clear()
	for row, vals := range rows {
		for col, text := range vals {
			cell := tview.NewTableCell(text)
			cell.SetExpansion(1)

			if row == 0 {
				cell.SetSelectable(false)
			}

			this.data.SetCell(row, col, cell)
		}
	}

	this.data.ScrollToBeginning()
}
