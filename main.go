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

	d, err := docker.NewDocker(wg, ctx)
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		os.Exit(1)
	}

	machine, err := state.New(state.DefaultConfig())
	if err != nil {
		panic(err)
	}

	tui := tui.NewTui(machine)
	wg.Go(func() {
		if err = tui.Run(); err != nil {
			cancel()
		}
	})

	wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				tui.Stop()
				return
			case err := <-d.ErrCh:
				tui.Send(err)
			case summary := <-d.ContainersCh:
				tui.Send(summary)
			case networks := <-d.NetworksCh:
				tui.Send(networks)
			case images := <-d.ImagesCh:
				tui.Send(images)
			}
		}
	})

	go func() {
		<-sigCh
		cancel()
		if tui != nil {
			tui.Stop()
		}
	}()

	d.Init()
	wg.Wait()
}
