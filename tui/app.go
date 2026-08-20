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
	menu        *tview.Flex
	nav         *tview.Pages
	pages       *tview.Pages
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
	this.wg.Go(func() {
		for {
			select {
			case <-this.ctx.Done():
				this.Stop()
				return
			case result := <-this.docker.Containers():
				this.containers.Update(result.Items)
			case result := <-this.docker.Images():
				this.images.Update(result.Items)
			case result := <-this.docker.Networks():
				this.networks.Update(result.Items)
			case err := <-this.docker.ErrCh:
				this.errors.Update(err)
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
	key := Utils.tcellEvent2Str(event)

	if state, ok := this.sm.Next(key); ok {
		switch state.Id {
		case "quit":
			this.cancel()
			return nil
		case "containers":
			this.containers.Focus()
		case "images":
			this.images.Focus()
		case "errors":
			this.errors.Focus()
		default:
			this.app.SetFocus(this.root)
		}

		this.pages.SwitchToPage(state.Id)
		this.nav.SwitchToPage(state.Id)

		return nil
	}

	return event
}
