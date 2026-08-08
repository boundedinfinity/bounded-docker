package main

import (
	"fmt"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/rivo/tview"
)

func newSummariesInfo(tui *tui, errs chan error) summariesInfo {
	this := summariesInfo{
		tui:   tui,
		items: make([]container.Summary, 0),
		route: listThing[client.ContainerListResult]{
			ch:   make(chan client.ContainerListResult),
			errs: errs,
			getfn: func() (client.ContainerListResult, error) {
				return tui.api.ContainerList(tui.ctx, client.ContainerListOptions{All: true})
			},
		},
	}

	this.create()
	this.draw()
	return this
}

type summariesInfo struct {
	tui   *tui
	items []container.Summary
	table *tview.Table
	route listThing[client.ContainerListResult]
}

func (this *summariesInfo) handle(result client.ContainerListResult) bool {
	if result.Items == nil || this.table == nil || this.tui.app == nil {
		return false
	}

	this.items = result.Items
	this.draw()
	return true
}

func (this *summariesInfo) draw() {
	this.table.Clear()
	headers := []string{"ID", "NAMES", "IMAGE", "STATUS"}
	this.table.SetTitle(fmt.Sprintf(" [ containers:%d ]", len(this.items)))

	type options struct {
		trunc bool
	}

	normal := func(text string, options options) string {
		text = strings.TrimSpace(text)

		if options.trunc && len(text) > this.tui.options.cellWidth {
			return text[:this.tui.options.cellWidth-3] + "..."
		}

		return text
	}

	var rows [][]string
	rows = append(rows, headers)

	for _, summary := range this.items {
		rows = append(rows, []string{
			normal(summary.ID, options{trunc: true}),
			normal(strings.Join(summary.Names, ", "), options{trunc: true}),
			normal(summary.Image, options{trunc: true}),
			normal(summary.Status, options{trunc: false}),
		})
	}

	for row, cols := range rows {
		for col, text := range cols {
			if col > 0 {
				text = this.tui.cellPadding(text)
			}

			cell := tview.NewTableCell(text)

			if row == 0 {
				cell.SetTextColor(this.tui.options.headerCorlor)
			} else {
				cell.SetTextColor(this.tui.options.foregroundColor)
			}

			this.table.SetCell(row, col, cell)
		}
	}
}

func (this *summariesInfo) create() {
	this.table = tview.NewTable()
	this.table.
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true).
		SetBackgroundColor(this.tui.options.backgroundColor).
		SetTitleColor(this.tui.options.titleColor)

	this.table.SetBorderColor(this.tui.options.foregroundColor)
}
