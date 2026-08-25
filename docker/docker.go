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
	ContainerLogs    *caller1[moby.ContainerLogsResult, string, moby.ContainerLogsOptions]
	ContainerInspect *caller1[moby.ContainerInspectResult, string, moby.ContainerInspectOptions]
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
		this.ContainerLogs.Init()
		this.ContainerInspect.Init()
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
					switch event.Action {
					case "start", "restart", "unpause":
						this.Containers.Run(nil)
					case "stop", "kill", "pause", "die":
						this.Containers.Run(nil)
					}
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
	if this.Api != nil {
		this.Api.Close()
	}
}
