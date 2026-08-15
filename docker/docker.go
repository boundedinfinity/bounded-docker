package docker

import (
	"context"
	"sync"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
	"github.com/moby/moby/client"
	moby "github.com/moby/moby/client"
)

const (
	__USER_AGENT = "bounded-docker/1.0.0"
)

func NewDocker(wg *sync.WaitGroup, ctx context.Context) (*dockerSystem, error) {
	api, err := moby.New(
		moby.FromEnv,
		moby.WithUserAgent(__USER_AGENT),
	)

	if err != nil {
		return nil, err
	}

	return &dockerSystem{
		ErrCh:        make(chan error),
		ContainersCh: make(chan []container.Summary),
		ImagesCh:     make(chan []image.Summary),

		wg:  wg,
		api: api,
		ctx: ctx,
	}, nil
}

type dockerSystem struct {
	ContainersCh chan []container.Summary
	ImagesCh     chan []image.Summary
	ErrCh        chan error
	wg           *sync.WaitGroup
	api          *client.Client
	ctx          context.Context
}

func (this *dockerSystem) Init() {
	this.wg.Go(func() {
		result := this.api.Events(this.ctx, moby.EventsListOptions{})
		for {
			select {
			case <-this.ctx.Done():
				return
			case err := <-result.Err:
				go func() { this.ErrCh <- err }()
			case event := <-result.Messages:
				switch event.Type {
				case "container":
					this.ListSummaries()
				case "image":
					this.ListImages()
				}
			}
		}
	})

	this.ListSummaries()
	this.ListImages()
}

func (this *dockerSystem) ListSummaries() {
	this.wg.Go(func() {
		o := client.ContainerListOptions{
			All: true,
		}
		if result, err := this.api.ContainerList(this.ctx, o); err == nil {
			go func() { this.ContainersCh <- result.Items }()
		} else {
			go func() { this.ErrCh <- err }()
		}
	})
}

func (this *dockerSystem) ListImages() {
	this.wg.Go(func() {
		o := client.ImageListOptions{
			All: true,
		}
		if result, err := this.api.ImageList(this.ctx, o); err == nil {
			go func() { this.ImagesCh <- result.Items }()
		} else {
			go func() { this.ErrCh <- err }()
		}
	})
}
