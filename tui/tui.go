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
	key := Utils.Tview.tcellEvent2Str(event)

	if state, cmd, ok := this.sm.Next(key); ok {
		this.pages.SwitchToPage(state.Id)
		this.nav.SwitchToPage(state.Id)

		switch state.Id {
		case "quit":
			this.cancel()
			return nil
		case "containers":
			this.app.SetFocus(this.containers.Data())
		case "images":
			this.app.SetFocus(this.images.Data())
		case "errors":
			this.app.SetFocus(this.errors.Data())
		case "networks":
			this.app.SetFocus(this.networks.Data())
		case "container-details":
			if id, ok := this.containers.Id(); ok {
				this.setStatus("Container: %s", id)
			} else {
				fmt.Println("Container: <none>")
			}
		}

		if cmd != nil {
			switch cmd.Name {
			case "up", "down":
				return event
			}
		}

		return nil
	}

	return event
}
