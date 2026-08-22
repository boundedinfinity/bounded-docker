package docker

import (
	"context"
	"sync"

	moby "github.com/moby/moby/client"
)

const (
	__USER_AGENT = "bounded-docker/1.0.0"
)

type System struct {
	containers       *caller0[moby.ContainerListResult, moby.ContainerListOptions]
	containerInspect *caller1[moby.ContainerInspectResult, string, moby.ContainerInspectOptions]
	containerLogs    *caller1[moby.ContainerLogsResult, string, moby.ContainerLogsOptions]
	images           *caller0[moby.ImageListResult, moby.ImageListOptions]
	networks         *caller0[moby.NetworkListResult, moby.NetworkListOptions]
	ErrCh            chan error
	wg               *sync.WaitGroup
	api              *moby.Client
	ctx              context.Context
}

func (this *System) Containers() chan moby.ContainerListResult {
	return this.containers.Out()
}

func (this *System) Images() chan moby.ImageListResult {
	return this.images.Out()
}

func (this *System) Networks() chan moby.NetworkListResult {
	return this.networks.Out()
}

func (this *System) Run() {
	this.wg.Go(func() {
		this.images.loop()
		this.containers.loop()
		this.networks.loop()
		result := this.api.Events(this.ctx, moby.EventsListOptions{})

		for {
			select {
			case <-this.ctx.Done():
				this.Stop()
				return
			case err := <-result.Err:
				go func() { this.ErrCh <- err }()
			case event := <-result.Messages:
				switch event.Type {
				case "container":
					this.ListContainers()
				case "image":
					this.ListImages()
				case "network":
					this.ListNetworks()
				}
			}
		}
	})
}

func (this *System) Stop() {

}

func (this *System) ListContainers() {
	this.containers.Run(nil)
}

func (this *System) ListImages() {
	this.images.Run(nil)
}

func (this *System) ListNetworks() {
	this.networks.Run(nil)
}

type logContext struct {
	result moby.ContainerLogsResult
	cancel context.CancelFunc
}

func (this *System) GetLogs(id string) {
	// this.wg.Go(func() {
	// 	o := client.ContainerLogsOptions{
	// 		ShowStdout: true,
	// 		ShowStderr: true,
	// 		Follow:     true,
	// 	}

	// 	ctx, cancel := context.WithCancel(this.ctx)

	// 	if result, err := this.api.ContainerLogs(ctx, id, o); err == nil {
	// 		// go func() { this.lo <- result.Items }()
	// 	} else {
	// 		go func() { this.ErrCh <- err }()
	// 	}
	// })
}
