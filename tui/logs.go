package tui

import (
	"context"
	"errors"
	"io"

	"github.com/moby/moby/api/pkg/stdcopy"
	moby "github.com/moby/moby/client"
	"github.com/rivo/tview"
)

func newLogs(tui *tui, id string, result moby.ContainerLogsResult) *logs {
	logs := &logs{tui: tui, containerId: id, result: result}
	logs.ctx, logs.cancel = context.WithCancel(tui.ctx)

	logs.view = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	logs.view.SetBorder(true).SetTitle(" [ Logs: " + id + " ] ")

	return logs
}

type logs struct {
	tui         *tui
	view        *tview.TextView
	containerId string
	result      moby.ContainerLogsResult
	ctx         context.Context
	cancel      context.CancelFunc
}

func (this *logs) Data() *tview.TextView {
	return this.view
}

func (this *logs) Start() {
	this.tui.wg.Go(func() {
		<-this.ctx.Done()
		this.result.Close()
	})

	this.tui.wg.Go(func() {
		defer this.cancel()

		w := &logWriter{tui: this.tui, view: this.view, ctx: this.ctx}

		// StdCopy blocks until the stream ends, forwarding each chunk as it arrives.
		if _, err := stdcopy.StdCopy(w, w, this.result); err != nil {
			if this.ctx.Err() == nil && !errors.Is(err, io.EOF) {
				select {
				case this.tui.docker.ErrCh <- err:
				case <-this.ctx.Done():
				}
			}
		}
	})
}

func (this *logs) Stop() {
	if this.cancel != nil {
		this.cancel()
	}
}

type logWriter struct {
	tui  *tui
	view *tview.TextView
	ctx  context.Context
}

var errLogsClosed = errors.New("log stream closed")

func (this *logWriter) Write(p []byte) (int, error) {
	if this.ctx.Err() != nil {
		return 0, errLogsClosed
	}

	// tview primitives must only be mutated from the event loop.
	chunk := append([]byte(nil), p...)
	done := make(chan struct{})

	// Not tracked by the WaitGroup: QueueUpdateDraw never returns once the app has stopped.
	go func() {
		defer close(done)
		this.tui.app.QueueUpdateDraw(func() {
			this.view.Write(chunk)
			this.view.ScrollToEnd()
		})
	}()

	select {
	case <-done:
		return len(p), nil
	case <-this.ctx.Done():
		return 0, errLogsClosed
	}
}
