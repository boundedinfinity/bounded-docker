package docker

import (
	"context"
	"sync"
)

func newCaller[T any, O any](wg *sync.WaitGroup, ctx context.Context, errCh chan error, options O, runFn func(context.Context, O) (T, error)) *caller[T, O] {
	return &caller[T, O]{
		wg:        wg,
		ctx:       ctx,
		options:   options,
		optionsCh: make(chan O),
		runFn:     runFn,
		runCh:     make(chan struct{}),
		outCh:     make(chan T),
		errCh:     errCh,
	}
}

type caller[T any, O any] struct {
	wg        *sync.WaitGroup
	ctx       context.Context
	options   O
	optionsCh chan O
	runFn     func(context.Context, O) (T, error)
	runCh     chan struct{}
	outCh     chan T
	errCh     chan error
}

func (this *caller[T, O]) Init() {
	this.wg.Go(func() {
		for {
			select {
			case <-this.ctx.Done():
				return
			case <-this.runCh:
				if this.runFn != nil {
					if result, err := this.runFn(this.ctx, this.options); err == nil {
						this.outCh <- result
					} else {
						this.errCh <- err
					}
				}
			case options := <-this.optionsCh:
				this.options = options
			}
		}
	})

	this.Run()
}

func (this *caller[T, O]) Stop() {
}

func (this *caller[T, O]) Out() chan T {
	return this.outCh
}

func (this *caller[T, O]) Run() {
	this.wg.Go(func() { this.runCh <- struct{}{} })
}

func (this *caller[T, O]) Send(result T) {
	this.wg.Go(func() { this.outCh <- result })
}
