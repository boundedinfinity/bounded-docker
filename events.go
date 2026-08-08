package main

import (
	"fmt"

	"github.com/boundedinfinity/docker-tui/docker"
	"github.com/moby/moby/api/types/events"
)

func (this *tui) eventHandle(message events.Message) bool {
	switch message.Type {
	case events.Type(docker.EventTypes.Container):
		this.summaries.route.get()
		return true
	default:
		this.errSend(fmt.Errorf("unhandled event: %s", message.Type))
	}

	return false
}
