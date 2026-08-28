package tui

import (
	"context"
	"fmt"
	"strings"
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
		wg:     wg,
		containerLogsOptions: moby.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Tail:       "500",
		},
	}

	tui.newMain()
	tui.header = tui.newHeader()
	tui.nav = tui.newNavigation()

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.header, 0, 1, false).
		AddItem(tui.pages, 0, 5, false).
		AddItem(tui.nav, 0, 1, false)

	tui.app.SetRoot(layout, true).EnableMouse(true)
	tui.app.SetInputCapture(tui.handleInput)
	tui.app.SetBeforeDrawFunc(tui.handleRedraw)
	tui.app.SetFocus(tui.pages)

	return tui
}

// newWelcome centers a left-aligned block of text within the available space.
func newWelcome(lines ...string) tview.Primitive {
	view := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(strings.Join(lines, "\n"))

	width := 0

	for _, line := range lines {
		if w := tview.TaggedStringWidth(line); w > width {
			width = w
		}
	}

	block := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(view, len(lines), 0, false).
		AddItem(nil, 0, 1, false)

	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(block, width, 0, false).
		AddItem(nil, 0, 1, false)
}

func (this *tui) newMain() {
	this.root = newWelcome(
		"Use the navigation items below to view:",
		"",
		"  - Containers",
		"  - Images",
		"  - Networks",
		"  - Errors",
	)

	this.containers = newInfoList(this, "Containers", moby.ContainerListOptions{All: true},
		Utils.Container.summary2rows, Utils.Container.id)
	this.images = newInfoList(this, "Images", moby.ImageListOptions{},
		Utils.Image.summary2rows, Utils.Image.id)
	this.networks = newInfoList(this, "Networks", moby.NetworkListOptions{},
		Utils.Network.summary2rows, Utils.Network.id)
	this.errors = newInfoList(this, "Errors", int(0), Utils.Error.error2rows, nil)

	this.pages = tview.NewPages().
		AddPage("root", this.root, true, true).
		AddPage("container.list", this.containers.Data(), true, false).
		AddPage("image.list", this.images.Data(), true, false).
		AddPage("network.list", this.networks.Data(), true, false).
		AddPage("errors", this.errors.Data(), true, false).
		SwitchToPage(this.sm.StartState().Id)
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

func (this *tui) setStatusf(format string, a ...any) {
	this.queueDraw(func() {
		this.status.Clear()
		fmt.Fprintf(this.status, format, a...)
	})
}

func (this *tui) setToastf(format string, a ...any) {
	this.queueDraw(func() {
		this.toast.Clear()
		fmt.Fprintf(this.toast, format, a...)
	})
}

func (this *tui) newHeader() *tview.Flex {
	this.status = tview.NewTextView().SetText(this.sm.Current.Id)
	this.toast = tview.NewTextView().SetText("toast")

	left := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(this.status, 0, 1, false).
		AddItem(this.toast, 0, 1, false)

	right := tview.NewFlex().
		SetDirection(tview.FlexRow)

	items := []Info{this.containers, this.images, this.networks, this.errors}
	for _, info := range items {
		right.AddItem(info.Header(), 0, 1, false)
	}

	header := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(left, 0, 1, false).
		AddItem(tview.NewTextView(), 0, 1, false).
		AddItem(right, 0, 1, false)

	header.SetBorder(true).SetTitle(" [ Bounded Docker ] ")

	return header
}
