package tui

import (
	"context"
	"fmt"
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
		status: tview.NewTextView().SetText("Welcome"),
		wg:     wg,
	}

	tui.containers = newInfo(tui, "Containers", Utils.Container.containerTitles(), Utils.Container.container2Row, Utils.Container.containerId)
	tui.images = newInfo(tui, "Images", Utils.Image.imageTitles(), Utils.Image.image2Row, Utils.Image.id)
	tui.networks = newInfo(tui, "Networks", Utils.Network.networkTitles(), Utils.Network.network2Row, Utils.Network.id)
	tui.errors = newInfo(tui, "Errors", Utils.Error.errorTitles(), Utils.Error.error2Row, nil)

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

	tui.menu = tui.newMenu(tui.status, []Info{tui.containers, tui.images, tui.networks, tui.errors})
	tui.nav = tui.newNavigation()

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.menu, 0, 1, false).
		AddItem(tui.pages, 0, 5, false).
		AddItem(tui.nav, 0, 1, false)

	tui.app.SetRoot(layout, true).EnableMouse(true)
	tui.app.SetInputCapture(tui.handleInput)
	tui.app.SetBeforeDrawFunc(tui.handleRedraw)
	tui.app.SetFocus(tui.pages)

	return tui
}

func (this *tui) newNavigation() *tview.Pages {
	pages := tview.NewPages()
	pages.SetBorder(true).SetTitle(" [ Navigation ]")

	for _, state := range this.sm.States() {
		pages.AddPage(state.Id, Utils.Tview.state2Table(state), true, false)
	}

	pages.SwitchToPage(this.sm.Current.Id)

	return pages
}

func (this *tui) setStatus(format string, a ...any) {
	this.wg.Go(func() {
		this.app.QueueUpdateDraw(func() {
			this.status.Clear()
			fmt.Fprintf(this.status, format, a...)
		})
	})
}

func (this *tui) newMenu(status tview.Primitive, items []Info) *tview.Flex {
	menu := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewTextView(), 0, 1, false).
		AddItem(status, 0, 2, false)

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
