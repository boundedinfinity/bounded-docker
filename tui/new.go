package tui

import (
	"context"
	"fmt"
	"sync"

	"github.com/boundedinfinity/docker-tui/docker"
	"github.com/boundedinfinity/docker-tui/state"
	moby "github.com/moby/moby/client"
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
		containerLogsOptions: moby.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Tail:       "500",
		},
	}

	tui.containers = newInfo(tui, "Containers", Utils.Container.titles(),
		moby.ContainerListOptions{All: true},
		Utils.Container.summary2rows, Utils.Container.id)
	tui.images = newInfo(tui, "Images", Utils.Image.titles(),
		moby.ImageListOptions{},
		Utils.Image.summary2row, Utils.Image.id)
	tui.networks = newInfo(tui, "Networks", Utils.Network.titles(),
		moby.NetworkListOptions{},
		Utils.Network.summary2rows, Utils.Network.id)
	tui.errors = newInfo(tui, "Errors", Utils.Error.errorTitles(), int(0), Utils.Error.error2Row, nil)

	tui.pages = tview.NewPages().
		AddPage("root", tui.root, true, true).
		AddPage("container.list", tui.containers.Data(), true, false).
		AddPage("image.list", tui.images.Data(), true, false).
		AddPage("network.list", tui.networks.Data(), true, false).
		AddPage("errors", tui.errors.Data(), true, false).
		SwitchToPage(tui.sm.StartState().Id)

	tui.menu = tui.newMenu(tui.status, []Info{tui.containers, tui.images, tui.networks, tui.errors})
	tui.nav = tui.newNavigation()

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
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
	pages.SetBorder(true).SetTitle(" [ Navigation ] ")

	for _, state := range this.sm.States() {
		pages.AddPage(state.Id, Utils.Tview.state2Table(state), true, false)
	}

	pages.SwitchToPage(this.sm.Current.Id)

	return pages
}

func (this *tui) setStatus(text string) {
	this.queueDraw(func() {
		this.status.Clear()
		fmt.Fprint(this.status, text)
	})
}

func (this *tui) setStatusf(format string, a ...any) {
	this.queueDraw(func() {
		this.status.Clear()
		fmt.Fprintf(this.status, format, a...)
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
