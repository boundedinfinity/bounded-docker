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
	Containers       *caller0[moby.ContainerListResult, moby.ContainerListOptions]
	containerInspect *caller1[moby.ContainerInspectResult, string, moby.ContainerInspectOptions]
	containerLogs    *caller1[moby.ContainerLogsResult, string, moby.ContainerLogsOptions]
	Images           *caller0[moby.ImageListResult, moby.ImageListOptions]
	Networks         *caller0[moby.NetworkListResult, moby.NetworkListOptions]
	ErrCh            chan error
	wg               *sync.WaitGroup
	Api              *moby.Client
	ctx              context.Context
}

func (this *System) Run() {
	this.wg.Go(func() {
		this.Images.loop()
		this.Containers.loop()
		this.Networks.loop()
		result := this.Api.Events(this.ctx, moby.EventsListOptions{})

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
					this.Containers.Run(nil)
				case "image":
					this.Images.Run(nil)
				case "network":
					this.Networks.Run(nil)
				}
			}
		}
	})
}

func (this *System) Stop() {

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
