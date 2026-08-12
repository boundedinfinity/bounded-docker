package docker

import (
	"context"
	"sync"

	"github.com/moby/moby/api/types/container"
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
		ErrCh:     make(chan error),
		SummaryCh: make(chan []container.Summary),
		wg:        wg,
		api:       api,
		ctx:       ctx,
	}, nil
}

type dockerSystem struct {
	SummaryCh chan []container.Summary
	ErrCh     chan error
	wg        *sync.WaitGroup
	api       *client.Client
	ctx       context.Context
}

func (this *dockerSystem) GetSummary() {
	this.wg.Go(func() {
		o := client.ContainerListOptions{
			All: true,
		}
		if result, err := this.api.ContainerList(this.ctx, o); err == nil {
			go func() { this.SummaryCh <- result.Items }()
		} else {
			go func() { this.ErrCh <- err }()
		}
	})
}
