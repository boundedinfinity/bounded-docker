package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/boundedinfinity/docker-tui/docker"
	"github.com/boundedinfinity/docker-tui/state"
	"github.com/boundedinfinity/docker-tui/tui"
)

func main() {
	var err error
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	d, err := docker.New(wg, ctx)
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		os.Exit(1)
	}

	wg.Go(func() {
		d.Run()
	})

	machine, err := state.New(state.DefaultConfig())
	if err != nil {
		panic(err)
	}

	tui := tui.New(wg, ctx, cancel, machine, d)
	wg.Go(func() {
		defer cancel()
		err = tui.Run()
	})

	go func() {
		<-sigCh
		cancel()
	}()

	wg.Wait()
}
