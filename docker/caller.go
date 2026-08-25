package docker

import (
	"context"
	"sync"
)

func newCaller0[T any, O any](wg *sync.WaitGroup, ctx context.Context, errCh chan error, runFn func(context.Context, O) (T, error)) *caller0[T, O] {
	var options O
	return &caller0[T, O]{
		wg:      wg,
		ctx:     ctx,
		options: options,
		runFn:   runFn,
		runCh:   make(chan *O),
		outCh:   make(chan T),
		errCh:   errCh,
	}
}

type caller0[T any, O any] struct {
	wg      *sync.WaitGroup
	ctx     context.Context
	options O
	runFn   func(context.Context, O) (T, error)
	runCh   chan *O
	outCh   chan T
	errCh   chan error
}

func (this *caller0[T, O]) loop() {
	this.wg.Go(func() {
		for {
			select {
			case <-this.ctx.Done():
				return
			case options := <-this.runCh:
				if options != nil {
					this.options = *options
				}

				if this.runFn != nil {
					if result, err := this.runFn(this.ctx, this.options); err == nil {
						this.outCh <- result
					} else {
						this.errCh <- err
					}
				}
			}
		}
	})
}

func (this *caller0[T, O]) Stop() {
}

func (this *caller0[T, O]) Out() chan T {
	return this.outCh
}

func (this *caller0[T, O]) Run(options *O) {
	this.wg.Go(func() { this.runCh <- options })
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type args1[A any, O any] struct {
	arg1    A
	options *O
}

func newCaller1[T any, A any, O any](wg *sync.WaitGroup, ctx context.Context, errCh chan error, runFn func(context.Context, A, O) (T, error)) *caller1[T, A, O] {
	var options O
	return &caller1[T, A, O]{
		wg:      wg,
		ctx:     ctx,
		options: options,
		runFn:   runFn,
		runCh:   make(chan args1[A, O]),
		outCh:   make(chan T),
		errCh:   errCh,
	}
}

type caller1[T any, A any, O any] struct {
	wg      *sync.WaitGroup
	ctx     context.Context
	options O
	runFn   func(context.Context, A, O) (T, error)
	runCh   chan args1[A, O]
	outCh   chan T
	errCh   chan error
}

func (this *caller1[T, A, O]) Init() {
	this.wg.Go(func() {
		for {
			select {
			case <-this.ctx.Done():
				return
			case args := <-this.runCh:
				if args.options != nil {
					this.options = *args.options
				}

				if this.runFn != nil {
					if result, err := this.runFn(this.ctx, args.arg1, this.options); err == nil {
						this.outCh <- result
					} else {
						this.errCh <- err
					}
				}
			}
		}
	})
}

func (this *caller1[T, A, O]) Stop() {
}

func (this *caller1[T, A, O]) Out() chan T {
	return this.outCh
}

func (this *caller1[T, A, O]) Run(arg1 A, options *O) {
	this.wg.Go(func() { this.runCh <- args1[A, O]{arg1: arg1, options: options} })
}
