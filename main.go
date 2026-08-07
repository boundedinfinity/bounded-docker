package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/boundedinfinity/docker-tui/docker"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/events"
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
	tui.runUpdateContainerList()
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
		api:          api,
		wg:           &sync.WaitGroup{},
		errCh:        make(chan error),
		itemsCh:      make(chan []container.Summary),
		containerMap: make(map[string]container.Summary),
		eventsCh:     make(chan events.Message),
		cellPadding:  strings.Repeat(" ", 4),
		options:      &tuiOptions{cellPadding: 1},
		optionsCh:    make(chan tuiOptions),
	}

	this.ctx, this.cancel = context.WithCancel(context.Background())

	this.createApp()
	this.createFlexBox()
	this.createTable()
	this.createErrorBox()
	this.flex.AddItem(this.table, 0, 1, false)
	this.flex.AddItem(this.errors, 0, 1, false)

	return this
}

type tui struct {
	ctx          context.Context
	cancel       context.CancelFunc
	api          *client.Client
	wg           *sync.WaitGroup
	app          *tview.Application
	flex         *tview.Flex
	table        *tview.Table
	errors       *tview.TextArea
	cellPadding  string
	headers      []string
	errCh        chan error
	itemsCh      chan []container.Summary
	containerMap map[string]container.Summary
	eventsCh     chan events.Message
	options      *tuiOptions
	optionsCh    chan tuiOptions
}

type tuiOptions struct {
	cellPadding int
}

func (this *tui) close() {
	this.app.Stop()
	this.api.Close()
	close(this.errCh)
	close(this.itemsCh)
}

func (this *tui) ui() {
	this.wg.Go(func() {
		if err := this.app.SetRoot(this.flex, true).SetFocus(this.table).Run(); err != nil {
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
			select {
			case <-this.ctx.Done():
				this.close()
				return
			case err := <-this.errCh:
				this.handleErr(err)
			case message := <-events.Messages:
				this.handleMessage(message)
			case err := <-events.Err:
				this.sendErr(err)
			case options := <-this.optionsCh:
				this.handleOptions(options)
			case items := <-this.itemsCh:
				this.updateMap(items)
				this.handlUpdateContainerList()
			}
		}
	})

	return gerr
}

func (this *tui) sendOptions(options tuiOptions) {
	go func() { this.optionsCh <- options }()
}

func (this *tui) handleOptions(options tuiOptions) {
	this.cellPadding = strings.Repeat(" ", options.cellPadding)
	this.app.ForceDraw()
}

func (this *tui) handleMessage(message events.Message) {
	switch message.Type {
	case events.Type(docker.EventTypes.Container):
		this.runUpdateContainerList()
	default:
		this.sendErr(fmt.Errorf("unhandled event type: %s", message.Type))
	}
}

func (this *tui) updateMap(items []container.Summary) {
	clear(this.containerMap)

	for _, item := range items {
		if _, ok := this.containerMap[item.ID]; !ok {
			this.containerMap[item.ID] = item
		} else {
			this.sendErr(fmt.Errorf("already seen container: %s", item.ID))
		}
	}
}

func (this *tui) sendErr(err error) {
	if err != nil {
		go func() { this.errCh <- err }()
	}
}

func (this *tui) handleErr(err error) {
	if err == nil || this.errors == nil || this.app == nil {
		return
	}

	text := strings.Join([]string{this.errors.GetText(), err.Error()}, "\n")
	this.errors.SetText(text, false)
	this.app.ForceDraw()
}

func (this *tui) clearErr() {
	if this.errors == nil || this.app == nil {
		return
	}

	this.errors.SetText("", false)
	this.app.ForceDraw()
}

func (this *tui) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	// bounded-docker.application.quit
	if event.Key() == tcell.KeyEscape || event.Rune() == 'q' || event.Rune() == 'Q' {
		this.cancel()
		return nil
	}

	// bounded-docker.application.errors.clear
	if event.Rune() == 'c' || event.Rune() == 'C' {
		this.clearErr()
		return nil
	}

	if event.Key() == tcell.KeyEnter {
		// row, col := this.table.GetSelection()
		// fmt.Println(row, col)
	}

	if event.Rune() == '+' || event.Rune() == '=' {
		current := *this.options
		current.cellPadding += 1
		this.sendOptions(current)
		return nil
	}

	if event.Rune() == '-' || event.Rune() == '_' {
		current := *this.options
		current.cellPadding -= 1
		this.sendOptions(current)
		return nil
	}

	return event
}

func (this *tui) createApp() {
	this.app = tview.NewApplication()
	this.app.SetInputCapture(this.handleKeyEvent)
}

func (this *tui) createFlexBox() {
	this.flex = tview.NewFlex().
		SetDirection(tview.FlexRow)
}

func (this *tui) createErrorBox() {
	this.errors = tview.NewTextArea()
	this.errors.SetTitle("Errors")
	this.errors.SetBorder(true)
}

func (this *tui) createTable() {
	this.headers = []string{"ID", "Names", "Image", "Status"}
	cols := len(this.headers) - 1

	this.table = tview.NewTable().
		SetBorders(true).
		SetFixed(1, cols).
		SetSelectable(true, false).
		SetDoneFunc(func(key tcell.Key) {
			// if key == tcell.KeyEnter {
			// 	this.table.SetSelectable(true, true)
			// }
		})

	this.table.SetTitle("Containers")

	for col := range this.headers {
		this.table.SetCell(0, col, createHeaderCell(this.headers[col]))
	}

	this.table.Select(0, 0)
}

func createHeaderCell(text string) *tview.TableCell {
	return tview.NewTableCell(text).
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
}

func (this *tui) createContainerCell(text string, id string) *tview.TableCell {
	text = this.cellPadding + text + this.cellPadding

	cell := tview.NewTableCell(text).
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignLeft).
		SetSelectable(true).
		SetReference(id)

	return cell
}

func (this *tui) runUpdateContainerList() {
	this.wg.Go(func() {
		if containerList, err := this.api.ContainerList(this.ctx, client.ContainerListOptions{All: true}); err == nil {
			go func() { this.itemsCh <- containerList.Items }()
		} else {
			this.sendErr(err)
		}
	})
}

func (this *tui) summary2Row(summary container.Summary) []string {
	type options struct {
		trunc bool
	}

	normal := func(text string, options options) string {
		text = strings.TrimSpace(text)

		if options.trunc && len(text) > 12 {
			return text[:12]
		}

		return text
	}

	return []string{
		normal(summary.ID, options{trunc: true}),
		normal(strings.Join(summary.Names, ", "), options{trunc: true}),
		normal(summary.Image, options{trunc: true}),
		normal(summary.Status, options{trunc: false}),
	}
}

func (this *tui) handlUpdateContainerList() {
	row := 0
	for _, item := range this.containerMap {
		row += 1
		texts := this.summary2Row(item)
		for col := range texts {
			cell := this.createContainerCell(texts[col], item.ID)
			this.table.SetCell(row, col, cell)
		}
	}

	this.app.ForceDraw()
}
