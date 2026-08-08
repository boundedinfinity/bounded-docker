package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/api/types/system"
	"github.com/moby/moby/client"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	tui := newTui()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	go func() {
		<-signalCh
		tui.cancel()
	}()

	tui.loop()
	tui.ui()
	tui.summaries.route.get()
	tui.dockerInfoSend()
	tui.wg.Wait()
}

func newTui() *tui {
	api, err := client.New(
		client.FromEnv,
		client.WithUserAgent("bounded-docker/1.0.0"),
	)
	if err != nil {
		panic(err)
	}

	this := &tui{
		api: api,
		wg:  &sync.WaitGroup{},
		errors: errorsInfo{
			ch:    make(chan error),
			items: make([]error, 0),
		},
		dockerInfo: dockerInfo{
			ch:   make(chan system.Info),
			info: system.Info{},
		},
		eventsCh:  make(chan events.Message),
		optionsCh: make(chan tuiOptions),
		options: tuiOptions{
			cellPadding:     10,
			cellWidth:       50,
			backgroundColor: tcell.ColorWhite,
			foregroundColor: tcell.ColorBlack,
			titleColor:      tcell.ColorBlack,
			headerCorlor:    tcell.ColorBlack,
		},
	}

	this.ctx, this.cancel = context.WithCancel(context.Background())

	this.createApp()
	this.errCreate()

	this.summaries = newSummariesInfo(this, this.errors.ch)
	this.infoCreate()
	this.flex.AddItem(this.summaries.table, 0, 1, false)
	this.flex.AddItem(this.errors.table, 0, 1, false)

	this.errDraw()

	return this
}

type tui struct {
	ctx              context.Context
	cancel           context.CancelFunc
	api              *client.Client
	wg               *sync.WaitGroup
	app              *tview.Application
	flex             *tview.Flex
	errors           errorsInfo
	summaries        summariesInfo
	dockerInfo       dockerInfo
	eventsCh         chan events.Message
	optionsCh        chan tuiOptions
	options          tuiOptions
	containersTitlef string
}

func (this *tui) close() {
	this.app.Stop()
	this.api.Close()
}

func (this *tui) ui() {
	this.wg.Go(func() {
		if err := this.app.SetRoot(this.flex, true).SetFocus(this.summaries.table).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	})
}

func (this *tui) loop() error {
	var gerr error

	this.wg.Go(func() {
		events := this.api.Events(this.ctx, client.EventsListOptions{})

		for {
			redraw := false

			select {
			case <-this.ctx.Done():
				this.close()
				return
			case err := <-this.errors.ch:
				redraw = this.errHandle(err)
			case message := <-events.Messages:
				redraw = this.eventHandle(message)
			case err := <-events.Err:
				this.errSend(err)
			case options := <-this.optionsCh:
				redraw = this.handleOptions(options)
			case result := <-this.summaries.route.ch:
				redraw = this.summaries.handle(result)
			case info := <-this.dockerInfo.ch:
				redraw = this.dockerInfoHandle(info)
			}

			if redraw {
				this.app.ForceDraw()
			}
		}
	})

	return gerr
}

func (this *tui) createApp() {
	this.app = tview.NewApplication()
	this.app.SetInputCapture(this.keyHandle)
	this.flex = tview.NewFlex().SetDirection(tview.FlexRow)
	// this.flex.
	// 	// SetBorderPadding(1, 1, 2, 2).
	// 	SetBackgroundColor(tcell.ColorWhite)
}

func newThing[T any](getfn func() (T, error), errs chan error) *listThing[T] {
	return &listThing[T]{
		ch:    make(chan T),
		getfn: getfn,
		errs:  errs,
	}
}

type listThing[T any] struct {
	ch       chan T
	getfn    func() (T, error)
	errs     chan error
	handleFn func(T) bool
}

func (this *listThing[T]) get() {
	if this.ch == nil || this.getfn == nil {
		return
	}

	go func() {
		if result, err := this.getfn(); err == nil {
			this.ch <- result
		} else {
			this.errs <- err
		}
	}()
}
