package tui

import (
	"fmt"
	"strings"

	"github.com/boundedinfinity/docker-tui/state"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/gdamore/tcell/v2"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
	"github.com/moby/moby/api/types/network"
	"github.com/rivo/tview"
)

func NewTui(sm *state.Machine) *tui {
	app := tview.NewApplication()

	tui := &tui{
		app:  app,
		sm:   sm,
		root: tview.NewTextView().SetText("Welcome"),
	}

	tui.containers = newInfo(tui, "Containers", containerTitles(), container2Row)
	tui.images = newInfo(tui, "Images", imageTitles(), image2Row)
	tui.networks = newInfo(tui, "Networks", networkTitles(), network2Row)
	tui.errors = newInfo(tui, "Errors", errorTitles(), error2Row)

	tui.containers.Init()
	tui.images.Init()
	tui.networks.Init()
	tui.errors.Init()

	tui.pages = tview.NewPages().
		AddPage("root", tui.root, true, true).
		AddPage("containers", tui.containers.Data(), true, true).
		AddPage("images", tui.images.Data(), true, false).
		AddPage("networks", tui.networks.Data(), true, false).
		AddPage("errors", tui.errors.Data(), true, false).
		SwitchToPage(tui.sm.Start().Id)

	menu := tui.newMenu([]Info{tui.containers, tui.images, tui.networks, tui.errors})

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(menu, 0, 1, false).
		AddItem(tui.pages, 0, 5, false).
		AddItem(tui.newNavigation(), 0, 1, false)

	tui.app.SetRoot(layout, true).EnableMouse(true)
	tui.app.SetInputCapture(tui.handleInput)
	tui.app.SetBeforeDrawFunc(tui.handleRedraw)

	return tui
}

type tui struct {
	app         *tview.Application
	middle      *tview.Box
	root        tview.Primitive
	containers  Info
	images      Info
	networks    Info
	errors      Info
	sm          *state.Machine
	help        *tview.Flex
	pages       *tview.Pages
	screenWidth int
}

func (this *tui) Send(event any) {
	switch e := event.(type) {
	case container.Summary, []container.Summary:
		this.containers.Update(e)
	case image.Summary, []image.Summary:
		this.images.Update(e)
	case network.Summary, []network.Summary:
		this.networks.Update(e)
	case error, []error:
		this.errors.Update(e)
	}
}

var (
	ErrTui   = errorer.New("tui")
	errTuiFn = errorer.Func(ErrTui)
)

func (t *tui) Run() error {
	return t.app.Run()
}

func (t *tui) Stop() {
	if t.app != nil {
		t.app.Stop()
	}
}

func (this *tui) handleRedraw(screen tcell.Screen) bool {
	width, _ := screen.Size()
	this.screenWidth = width
	return false
}

func tcellEvent2Str(event *tcell.EventKey) string {
	key := string(event.Name())
	key = strings.Replace(key, "Rune[", "", 1)
	key = strings.Replace(key, "]", "", 1)
	key = strings.ToLower(key)

	return key
}

func (this *tui) handleInput(event *tcell.EventKey) *tcell.EventKey {
	key := tcellEvent2Str(event)

	if state, ok := this.sm.Next(key); ok {
		this.pages.SwitchToPage(state.Id)
		this.updateNavigation()

		switch state.Id {
		case "containers":
			this.containers.Focus()
		case "images":
			this.images.Focus()
		case "errors":
			this.errors.Focus()
		default:
			this.app.SetFocus(this.root)
		}

		return nil
	}

	return event
}

func (this *tui) newNavigation() tview.Primitive {
	this.help = tview.NewFlex().SetDirection(tview.FlexColumn)
	this.help.SetBorder(true).SetTitle(" [ Navigation ]")
	this.updateNavigation()

	return this.help
}

func (this *tui) updateNavigation() {
	this.help.Clear()

	this.help.AddItem(tview.NewTextView(), 0, 1, false)
	for _, navigation := range this.sm.Current.Transitions {
		var keys []string
		for _, key := range navigation.Keys {
			keys = append(keys, key.Name)
		}

		text := fmt.Sprintf("%s -> %s", strings.Join(keys, "/"), navigation.State.Name)
		view := tview.NewTextView().SetText(text)
		this.help.AddItem(view, 0, 1, false)
	}
	this.help.AddItem(tview.NewTextView(), 0, 1, false)
}

func (this *tui) newMenu(items []Info) tview.Primitive {
	menu := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewTextView(), 0, 1, false)

	for _, info := range items {
		menu.AddItem(info.Header(), 0, 1, false)
	}

	menu.
		AddItem(tview.NewTextView(), 0, 1, false).
		SetBorder(true).
		SetTitle(" [ Bounded Docker ] ")

	return menu
}
