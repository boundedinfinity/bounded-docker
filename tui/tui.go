package tui

import (
	"context"
	"fmt"
	"sync"

	"github.com/boundedinfinity/docker-tui/docker"
	"github.com/boundedinfinity/docker-tui/state"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/gdamore/tcell/v2"
	moby "github.com/moby/moby/client"
	"github.com/rivo/tview"
)

type tui struct {
	docker           *docker.System
	ctx              context.Context
	cancel           context.CancelFunc
	app              *tview.Application
	middle           *tview.Box
	root             tview.Primitive
	containers       Info
	containerOptions moby.ContainerListOptions
	images           Info
	imageOptions     moby.ImageListOptions
	networks         Info
	errors           Info
	sm               *state.Machine
	status           *tview.TextView
	menu             *tview.Flex
	nav              *tview.Pages
	pages            *tview.Pages
	current          string
	screenWidth      int
	logsCh           chan moby.ContainerLogsResult
	logsCtx          context.Context
	wg               *sync.WaitGroup
}

var (
	ErrTui   = errorer.New("tui")
	errTuiFn = errorer.Func(ErrTui)
)

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

			if result, err := this.docker.Api.ContainerList(this.ctx, this.containerOptions); err != nil {
				this.errors.Queue(err)
			} else {
				this.containers.Queue(result.Items)
			}
		})
	}

	if c, ok := this.sm.GetCommand("image.list.all"); ok {
		c.AddRunFunc(func(_ *state.State, _ *state.Command) {
			this.imageOptions.All = !this.imageOptions.All

			if result, err := this.docker.Api.ImageList(this.ctx, this.imageOptions); err != nil {
				this.errors.Queue(err)
			} else {
				this.images.Queue(result.Items)
			}
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
			case result := <-this.logsCh:
				defer result.Close()
				fmt.Printf("%v\n", result)
			}
		}
	})

	this.docker.Containers.Run(&moby.ContainerListOptions{All: true})
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
