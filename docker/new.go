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

	s := &System{
		ErrCh:      errCh,
		wg:         wg,
		Api:        api,
		ctx:        ctx,
		Containers: newCaller0(wg, ctx, errCh, api.ContainerList),
		Images:     newCaller0(wg, ctx, errCh, api.ImageList),
		Networks:   newCaller0(wg, ctx, errCh, api.NetworkList),
	}

	return s, nil
}
