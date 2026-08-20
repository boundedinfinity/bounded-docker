package docker

import (
	"context"
	"sync"

	moby "github.com/moby/moby/client"
)

func New(wg *sync.WaitGroup, ctx context.Context) (*System, error) {
	api, err := moby.New(
		moby.FromEnv,
		moby.WithUserAgent(__USER_AGENT),
	)

	if err != nil {
		return nil, err
	}

	errCh := make(chan error)

	return &System{
		ErrCh:      errCh,
		wg:         wg,
		api:        api,
		ctx:        ctx,
		containers: newCaller(wg, ctx, errCh, moby.ContainerListOptions{All: true}, api.ContainerList),
		images:     newCaller(wg, ctx, errCh, moby.ImageListOptions{All: true}, api.ImageList),
		networks:   newCaller(wg, ctx, errCh, moby.NetworkListOptions{}, api.NetworkList),
	}, nil
}
