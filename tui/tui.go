package tui

import (
	"context"
	"sync"

	"github.com/boundedinfinity/docker-tui/docker"
	"github.com/boundedinfinity/docker-tui/state"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/gdamore/tcell/v2"
	moby "github.com/moby/moby/client"
	"github.com/rivo/tview"
)

type tui struct {
	docker               *docker.System
	ctx                  context.Context
	cancel               context.CancelFunc
	app                  *tview.Application
	middle               *tview.Box
	root                 tview.Primitive
	containers           Info
	containerOptions     moby.ContainerListOptions
	containerLogs        *logs
	containerLogsOptions moby.ContainerLogsOptions
	images               Info
	imageOptions         moby.ImageListOptions
	networks             Info
	errors               Info
	sm                   *state.Machine
	status               *tview.TextView
	menu                 *tview.Flex
	nav                  *tview.Pages
	pages                *tview.Pages
	current              string
	screenWidth          int
	wg                   *sync.WaitGroup
}

var (
	ErrTui   = errorer.New("tui")
	errTuiFn = errorer.Func(ErrTui)
)

func (_ tui) createPid(containerId string) string {
	return "container.logs." + containerId
}

// openLogs replaces any running log stream with a fresh one for the given container.
func (this *tui) openLogs(containerId string) {
	if this.containerLogs != nil {
		this.containerLogs.Stop()
		this.pages.RemovePage(this.createPid(this.containerLogs.containerId))
		this.containerLogs = nil
	}

	result, err := this.docker.Api.ContainerLogs(this.ctx, containerId, this.containerLogsOptions)
	if err != nil {
		this.wg.Go(func() { this.docker.ErrCh <- err })
		return
	}

	pid := this.createPid(containerId)
	view := newLogs(this, containerId, result)
	this.containerLogs = view
	this.pages.AddPage(pid, view.Data(), true, false)
	this.pages.SwitchToPage(pid)
	view.Start()
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

	if c, ok := this.sm.GetCommand("container.logs.follow"); ok {
		c.AddRunFunc(func(_ *state.State, _ *state.Command) {
			this.containerLogsOptions.Follow = !this.containerLogsOptions.Follow
			if id, ok := this.containers.Id(); ok {
				this.docker.ContainerLogs.Run(id, &this.containerLogsOptions)
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

func (this *tui) handleInput(event *tcell.EventKey) *tcell.EventKey {
	key := Utils.Tview.tcellEvent2Str(event)
	state, _ := this.sm.Next(key)

	this.setStatus(state.Id)

	if this.current != state.Id {
		this.current = state.Id
		event = nil
	}

	return event
}
