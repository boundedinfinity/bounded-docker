package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type errorsInfo struct {
	ch    chan error
	items []error
	table *tview.Table
}

func (this *tui) errHandle(err error) bool {
	if err == nil || this.errors.table == nil || this.app == nil {
		return false
	}

	this.errors.items = append(this.errors.items, err)
	this.errDraw()
	return true
}

func (this *tui) errDraw() {
	this.errors.table.SetTitle(fmt.Sprintf(" [ errors:%d ]", len(this.errors.items)))
	this.errors.table.Clear()

	for row, err := range this.errors.items {
		cell := tview.NewTableCell(err.Error()).
			SetTextColor(this.options.foregroundColor)
		this.errors.table.SetCell(row, 0, cell)
	}
}

func (this *tui) errSend(err error) {
	if err != nil && this.errors.ch != nil {
		go func() { this.errors.ch <- err }()
	}
}

func (this *tui) errClear() {
	clear(this.errors.items)
	this.errDraw()
}

func (this *tui) errCreate() {
	this.errors.table = tview.NewTable()
	this.errors.table.
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true).
		SetBackgroundColor(this.options.backgroundColor).
		SetBorderColor(this.options.foregroundColor).
		SetTitleColor(this.options.titleColor)
}
