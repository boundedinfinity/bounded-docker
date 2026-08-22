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
	docker      *docker.System
	ctx         context.Context
	cancel      context.CancelFunc
	app         *tview.Application
	middle      *tview.Box
	root        tview.Primitive
	containers  Info
	images      Info
	networks    Info
	errors      Info
	sm          *state.Machine
	status      *tview.TextView
	menu        *tview.Flex
	nav         *tview.Pages
	pages       *tview.Pages
	current     string
	screenWidth int
	logsCh      chan moby.ContainerLogsResult
	logsCtx     context.Context
	wg          *sync.WaitGroup
}

var (
	ErrTui   = errorer.New("tui")
	errTuiFn = errorer.Func(ErrTui)
)

func (this *tui) Run() error {
	// queue := func(fn func()) {
	// 	this.wg.Go(func() {
	// 		this.app.QueueUpdate(fn)
	// 		this.app.Draw()
	// 	})
	// }

	this.wg.Go(func() {
		for {
			select {
			case <-this.ctx.Done():
				this.Stop()
				return
			case result := <-this.docker.Containers():
				this.containers.Queue(result.Items)
			case result := <-this.docker.Images():
				this.images.Queue(result.Items)
			case result := <-this.docker.Networks():
				this.networks.Queue(result.Items)
			case err := <-this.docker.ErrCh:
				this.errors.Queue(err)
			case result := <-this.logsCh:
				defer result.Close()
				fmt.Printf("%v\n", result)
			}
		}
	})

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
	state, cmd := this.sm.Next(key)

	if cmd != nil {
		switch cmd.Id {
		case "quit":
			this.cancel()
			return nil
		case "up", "down":
			return event
		}
	}

	if this.current != state.Id {
		this.current = state.Id
		this.pages.SwitchToPage(state.Id)
		this.nav.SwitchToPage(state.Id)
		// this.app.Draw()

		switch state.Id {
		case "container.list":
			this.docker.ListContainers()
		case "container.details":
			text := "Container"
			if id, ok := this.containers.Id(); ok {
				text += "[" + id + "]"
			}

			this.setStatus(text)
		}

		this.app.ForceDraw()
		return nil
	}

	return event
}
