package tui

import (
	"fmt"
	"strings"

	"github.com/boundedinfinity/docker-tui/state"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
	"github.com/rivo/tview"
)

func NewTui(sm *state.Machine) *tui {
	tui := &tui{
		sm:         sm,
		containers: createInfo("Containers", containerTitles(), container2Row),
		images:     createInfo("Images", imageTitles(), image2Row),
		errors:     createInfo("Errors", errorTitles(), error2Row),
	}

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.mkMenu(), 0, 1, false).
		AddItem(tui.containers.data, 0, 5, false).
		AddItem(tui.mkNavigation(), 0, 1, false)

	tui.app = tview.NewApplication().SetRoot(layout, true).EnableMouse(true)

	tui.containers.Update(createFakeContainers())
	tui.images.Update(createFakeImages())
	tui.errors.Update(createFakeErrors())

	return tui
}

type tui struct {
	app        *tview.Application
	middle     *tview.Box
	containers *Info[container.Summary]
	images     *Info[image.Summary]
	errors     *Info[error]
	sm         *state.Machine
	help       *tview.Flex
}

func (t *tui) Run() error {
	return t.app.Run()
}

func (this *tui) mkNavigation() tview.Primitive {
	this.help = tview.NewFlex().SetDirection(tview.FlexColumn)
	this.help.SetBorder(true).SetTitle(" [ Navigation ]")
	this.updateNavigation()

	return this.help
}

func (this *tui) updateNavigation() {
	this.help.Clear()
	var keys []string

	for _, navigation := range this.sm.Current.Navigations {
		for _, key := range navigation.Keys {
			keys = append(keys, key.Name)
		}

		text := fmt.Sprintf("%s -> %s", strings.Join(keys, "/"), navigation.Name)
		view := tview.NewTextView().SetText(text)
		this.help.AddItem(view, 0, 1, false)
	}
}

func (this *tui) mkMenu() tview.Primitive {
	menu := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewTextView(), 0, 1, false).
		AddItem(this.containers.header, 0, 1, false).
		AddItem(this.images.header, 0, 1, false).
		AddItem(this.errors.header, 0, 1, false).
		AddItem(tview.NewTextView(), 0, 1, false)
	menu.
		SetBorder(true).
		SetTitle(" [ Bounded Docker ] ")

	return menu
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

func createInfo[T any](title string, headers []string, fn func(i int, item T) []string) *Info[T] {
	info := &Info[T]{
		Title:        title,
		Items:        []T{},
		headers:      headers,
		item2RowFunc: fn,
	}

	info.header = tview.NewTextView()

	info.data = tview.NewTable().
		SetBorders(true).
		SetFixed(1, len(headers)).
		SetEvaluateAllRows(true)

	info.Update(nil)
	return info
}

type Info[T any] struct {
	Title        string
	Items        []T
	headers      []string
	header       *tview.TextView
	data         *tview.Table
	item2RowFunc func(int, T) []string
	colWidth     int
}

func (this *Info[T]) Update(items []T) {
	this.Items = items
	count := len(this.Items)

	this.header.Clear()
	fmt.Fprintf(this.header, "%s [%d]", this.Title, count)

	this.data.Clear()
	data := [][]string{this.headers}
	for i, item := range this.Items {
		data = append(data, this.item2RowFunc(i, item))
	}

	for row, vals := range data {
		for col, val := range vals {
			cell := tview.NewTableCell(val)
			this.data.SetCell(row, col, cell)
		}
	}
}
