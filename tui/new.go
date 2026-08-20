package tui

import (
	"context"
	"sync"

	"github.com/boundedinfinity/docker-tui/docker"
	"github.com/boundedinfinity/docker-tui/state"
	"github.com/rivo/tview"
)

func New(wg *sync.WaitGroup, ctx context.Context, cancel context.CancelFunc, sm *state.Machine, docker *docker.System) *tui {
	app := tview.NewApplication()

	tui := &tui{
		docker: docker,
		ctx:    ctx,
		app:    app,
		sm:     sm,
		cancel: cancel,
		root:   tview.NewTextView().SetText("Welcome"),
		wg:     wg,
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

	tui.menu = tui.newMenu([]Info{tui.containers, tui.images, tui.networks, tui.errors})
	tui.nav = tui.newNavigation()

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.menu, 0, 1, false).
		AddItem(tui.pages, 0, 5, false).
		AddItem(tui.nav, 0, 1, false)

	tui.app.SetRoot(layout, true).EnableMouse(true)
	tui.app.SetInputCapture(tui.handleInput)
	tui.app.SetBeforeDrawFunc(tui.handleRedraw)

	return tui
}

func (this *tui) newNavigation() *tview.Pages {
	pages := tview.NewPages()
	pages.SetBorder(true).SetTitle(" [ Navigation ]")

	for _, state := range this.sm.States() {
		pages.AddPage(state.Id, Utils.state2Table(state), true, false)
	}

	pages.SwitchToPage(this.sm.Current.Id)

	return pages
}

func (this *tui) newMenu(items []Info) *tview.Flex {
	menu := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewTextView(), 0, 1, false)
	menu.SetBorder(true).SetTitle(" [ Bounded Docker ] ")

	for _, info := range items {
		menu.AddItem(info.Header(), 0, 1, false)
	}

	menu.
		AddItem(tview.NewTextView(), 0, 1, false).
		SetBorder(true).
		SetTitle(" [ Bounded Docker ] ")

	return menu
}
