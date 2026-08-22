package state

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var Utils = utils{}

type utils struct {
	S sv
	K kv
	M mv
}

type sv struct{}

func (_ sv) norm(s string) string {
	return strings.TrimSpace(s)
}

type kv struct{}

func (_ kv) valid(s string) bool {
	return slices.Contains(validKeyCodes, s)
}

type mv struct{}

type stateKeyContext struct {
	stateId      string
	transitionId string
	commandId    string
	index        int
}

type stateKeyValidator map[string][]stateKeyContext

func (this stateKeyValidator) add(stateId, transitionId, commandId, key string, index int) {
	ctx := stateKeyContext{
		stateId:      stateId,
		transitionId: transitionId,
		commandId:    commandId,
		index:        index,
	}

	this[key] = append(this[key], ctx)
}

func (this stateKeyValidator) validate() error {
	var errs []error

	for key, ctxs := range this {
		if len(ctxs) > 1 {
			for _, ctx := range ctxs {
				item := fmt.Sprintf("[state:%s]", ctx.stateId)
				if ctx.transitionId != "" {
					item += fmt.Sprintf("[transition:%s]", ctx.transitionId)
				}

				if ctx.commandId != "" {
					item += fmt.Sprintf("[command:%s]", ctx.commandId)
				}

				item += fmt.Sprintf("[%s][%d]", key, ctx.index)
				errs = append(errs, errFn(item))
			}
		}
	}

	return errors.Join(errs...)
}
