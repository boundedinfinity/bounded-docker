package tui

import (
	"context"
	"sync"
	"time"

	"github.com/boundedinfinity/docker-tui/docker"
	"github.com/boundedinfinity/docker-tui/state"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/gdamore/tcell/v2"
	moby "github.com/moby/moby/client"
	"github.com/rivo/tview"
)

type tui struct {
	docker                  *docker.System
	ctx                     context.Context
	cancel                  context.CancelFunc
	app                     *tview.Application
	middle                  *tview.Box
	root                    tview.Primitive
	containers              Info
	containerOptions        moby.ContainerListOptions
	containerLogs           *logs
	containerLogsOptions    moby.ContainerLogsOptions
	containerInspect        *inspect
	containerInspectOptions moby.ContainerInspectOptions
	images                  Info
	imageOptions            moby.ImageListOptions
	networks                Info
	errors                  Info
	sm                      *state.Machine
	status                  *tview.TextView
	toast                   *tview.TextView
	header                  *tview.Flex
	nav                     *tview.Pages
	navTables               map[string]*tview.Table
	pages                   *tview.Pages
	current                 string
	screen                  tcell.Screen
	screenWidth             int
	wg                      *sync.WaitGroup
}

var (
	ErrTui   = errorer.New("tui")
	errTuiFn = errorer.Func(ErrTui)
)

func (_ tui) createPid(stateId, containerId string) string {
	return stateId + "." + containerId
}

// openLogs replaces any running log stream with a fresh one for the given container.
func (this *tui) openLogs(containerId string) {
	if this.containerLogs != nil {
		this.containerLogs.Stop()
		this.pages.RemovePage(this.createPid("container.logs", this.containerLogs.containerId))
		this.containerLogs = nil
	}

	result, err := this.docker.Api.ContainerLogs(this.ctx, containerId, this.containerLogsOptions)
	if err != nil {
		this.wg.Go(func() { this.docker.ErrCh <- err })
		return
	}

	pid := this.createPid("container.logs", containerId)
	view := newLogs(this, containerId, result)
	this.containerLogs = view
	this.pages.AddPage(pid, view.Data(), true, false)
	this.pages.SwitchToPage(pid)
	view.Start()
}

// openInspect shows an inspect view for the given container and requests its data.
func (this *tui) openInspect(containerId string) {
	if this.containerInspect != nil {
		this.pages.RemovePage(this.createPid("container.inspect", this.containerInspect.containerId))
		this.containerInspect = nil
	}

	pid := this.createPid("container.inspect", containerId)
	view := newInspect(this, containerId)
	this.containerInspect = view
	this.pages.AddPage(pid, view.Data(), true, false)
	this.pages.SwitchToPage(pid)

	this.docker.ContainerInspect.Run(containerId, &this.containerInspectOptions)
}

func (this *tui) Run() error {
	if c, ok := this.sm.GetCommand("quit"); ok {
		c.AddRunFunc(func(_ *state.State, _ *state.Command) {
			this.Stop()
		})
	}

	switchToPage := func(state *state.State) {
		this.pages.SwitchToPage(state.Id)
		this.nav.SwitchToPage(state.Id)
	}

	if s, ok := this.sm.GetState("root"); ok {
		s.AddEnterFunc(switchToPage)
	}

	if s, ok := this.sm.GetState("errors"); ok {
		s.AddEnterFunc(switchToPage)
	}

	if s, ok := this.sm.GetState("container.list"); ok {
		s.AddEnterFunc(switchToPage)
	}

	if s, ok := this.sm.GetState("image.list"); ok {
		s.AddEnterFunc(switchToPage)
	}

	if s, ok := this.sm.GetState("network.list"); ok {
		s.AddEnterFunc(switchToPage)
	}

	if c, ok := this.sm.GetCommand("container.list.all"); ok {
		c.AddRunFunc(func(_ *state.State, _ *state.Command) {
			this.containerOptions.All = !this.containerOptions.All
			this.docker.Containers.Run(&this.containerOptions)
		})
	}

	if s, ok := this.sm.GetState("container.logs"); ok {
		s.AddEnterFunc(func(state *state.State) {
			this.nav.SwitchToPage(state.Id)

			if cid, ok := this.containers.Id(); ok {
				this.openLogs(cid)
			}
		})
	}

	if c, ok := this.sm.GetCommand("container.logs.follow"); ok {
		c.AddRunFunc(func(_ *state.State, _ *state.Command) {
			this.containerLogsOptions.Follow = !this.containerLogsOptions.Follow

			if this.containerLogs != nil {
				this.openLogs(this.containerLogs.containerId)
			}
		})
	}

	if s, ok := this.sm.GetState("container.inspect"); ok {
		s.AddEnterFunc(func(state *state.State) {
			this.nav.SwitchToPage(state.Id)

			if cid, ok := this.containers.Id(); ok {
				this.openInspect(cid)
			}
		})
	}

	if c, ok := this.sm.GetCommand("container.inspect.size"); ok {
		c.AddRunFunc(func(_ *state.State, _ *state.Command) {
			this.containerInspectOptions.Size = !this.containerInspectOptions.Size

			if this.containerInspect != nil {
				this.docker.ContainerInspect.Run(this.containerInspect.containerId, &this.containerInspectOptions)
			}
		})
	}

	if c, ok := this.sm.GetCommand("container.inspect.copy"); ok {
		c.AddRunFunc(func(_ *state.State, _ *state.Command) {
			if this.containerInspect == nil {
				return
			}

			if value, ok := this.containerInspect.Value(); ok {
				this.copyToClipboard(value)
			}
		})
	}

	if c, ok := this.sm.GetCommand("image.list.all"); ok {
		c.AddRunFunc(func(_ *state.State, _ *state.Command) {
			this.imageOptions.All = !this.imageOptions.All
			this.docker.Images.Run(&this.imageOptions)
		})
	}

	this.wg.Go(func() {
		for {
			select {
			case <-this.ctx.Done():
				this.Stop()
				return
			case result := <-this.docker.Containers.Out():
				this.containers.Queue(result.Items)
			case result := <-this.docker.Images.Out():
				this.images.Queue(result.Items)
			case result := <-this.docker.Networks.Out():
				this.networks.Queue(result.Items)
			case result := <-this.docker.ContainerInspect.Out():
				if this.containerInspect != nil {
					this.containerInspect.Set(result)
				}
			case err := <-this.docker.ErrCh:
				this.errors.Queue(err)
			}
		}
	})

	this.docker.Containers.Run(nil)
	this.docker.Images.Run(nil)
	this.docker.Networks.Run(nil)
	return this.app.Run()
}

func (this *tui) Stop() {
	if this.cancel != nil {
		this.cancel()
	}

	if this.app != nil {
		this.app.Stop()
	}
}

// queueDraw applies fn on the event loop. It is deliberately not tracked by the
// WaitGroup because QueueUpdateDraw never returns once the app has stopped.
func (this *tui) queueDraw(fn func()) {
	if this.ctx.Err() != nil {
		return
	}

	go this.app.QueueUpdateDraw(fn)
}

func (this *tui) handleRedraw(screen tcell.Screen) bool {
	width, _ := screen.Size()
	this.screen = screen
	this.screenWidth = width
	return false
}

// copyToClipboard posts text to the system clipboard via the terminal's OSC 52
// support, which not every terminal honors.
func (this *tui) copyToClipboard(text string) {
	if this.screen == nil {
		this.setStatusf("clipboard unavailable")
		return
	}

	this.screen.SetClipboard([]byte(text))
	this.setToastf("copied %q", Utils.String.clamp(text, 50))
}

func (this *tui) handleInput(event *tcell.EventKey) *tcell.EventKey {
	key := Utils.Tview.tcellEvent2Str(event)
	state, _ := this.sm.Next(key)

	this.highlightKey(key)
	this.setStatusf("%s", state.Id)

	if this.current != state.Id {
		this.current = state.Id
		event = nil
	}

	return event
}

// highlightKey finds the given key code in the current navigation table and
// briefly changes its cell color to a bright contrasting color, then fades back
// to the normal color after 500ms.
func (this *tui) highlightKey(keyCode string) {
	table, ok := this.navTables[this.current]
	if !ok {
		return
	}

	highlightColor := tcell.ColorLimeGreen
	normalColor := tcell.ColorWhite
	normalKeyColor := tcell.ColorYellow

	this.queueDraw(func() {
		for r := 0; r < table.GetRowCount(); r++ {
			for c := 0; c < table.GetColumnCount(); c++ {
				cell := table.GetCell(r, c)
				if cell == nil {
					continue
				}
				if cell.Text == keyCode {
					cell.SetTextColor(highlightColor)
				}
			}
		}
	})

	// Schedule revert after 500ms
	go func() {
		select {
		case <-this.ctx.Done():
			return
		case <-time.After(500 * time.Millisecond):
		}

		this.queueDraw(func() {
			for r := 0; r < table.GetRowCount(); r++ {
				for c := 0; c < table.GetColumnCount(); c++ {
					cell := table.GetCell(r, c)
					if cell == nil {
						continue
					}
					if cell.Text == keyCode {
						if c%2 == 0 {
							cell.SetTextColor(normalColor)
						} else {
							cell.SetTextColor(normalKeyColor)
						}
					}
				}
			}
		})
	}()
}
