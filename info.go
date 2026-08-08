package main

import (
	"fmt"

	"github.com/moby/moby/api/types/system"
	"github.com/moby/moby/client"
	"github.com/rivo/tview"
)

type dockerInfo struct {
	ch    chan system.Info
	info  system.Info
	table *tview.Table
}

func (this *tui) dockerInfoSend() {
	this.wg.Go(func() {
		if result, err := this.api.Info(this.ctx, client.InfoOptions{}); err == nil {
			go func() { this.dockerInfo.ch <- result.Info }()
		} else {
			this.errSend(err)
		}
	})
}

func (this *tui) dockerInfoHandle(info system.Info) bool {
	if this.dockerInfo.table == nil || this.app == nil {
		return false
	}

	this.dockerInfo.info = info
	this.dockerInfoDraw()
	return true
}

func (this *tui) dockerInfoDraw() {
	this.dockerInfo.table.SetTitle(" Docker Info ")
	this.dockerInfo.table.Clear()

	info := this.dockerInfo.info
	row := 0
	for key, value := range map[string]interface{}{
		"ID":              info.ID,
		"Containers":      info.Containers,
		"Images":          info.Images,
		"Driver":          info.Driver,
		"Memory":          info.MemTotal,
		"OperatingSystem": info.OperatingSystem,
	} {
		cell := tview.NewTableCell(fmt.Sprintf("%s: %v", key, value)).
			SetTextColor(this.options.foregroundColor)
		this.dockerInfo.table.SetCell(row, 0, cell)
		row++
	}
}

func (this *tui) infoCreate() {
	this.dockerInfo.table = tview.NewTable()
	this.dockerInfo.table.
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true).
		SetBackgroundColor(this.options.backgroundColor).
		SetBorderColor(this.options.foregroundColor).
		SetTitleColor(this.options.titleColor)
}
