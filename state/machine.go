package state

import (
	"strings"

	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	ErrMachine = errorer.New("state machine")
	ErrConfig  = errorer.New("config")
	errFn      = errorer.Func(ErrMachine, ErrConfig)
)

// /////////////////////////////////////////////////////////////////////////////////////////////////

type Machine struct {
	Current     *State
	transitions []*State
	navigation  []*State
	keys        []*Key
	capture     strings.Builder
	selected    []string
}

func (this Machine) FindState(id string) (*State, bool) {
	for _, state := range this.transitions {
		if state.Id == id {
			return state, true
		}
	}

	return nil, false
}

func (this Machine) FindKey(code string) (*Key, bool) {
	for _, key := range this.keys {
		if key.Code == code {
			return key, true
		}
	}

	return nil, false
}

func (this Machine) Start() *State {
	return this.Current
}

func (this *Machine) Next(code string) (*State, bool) {
	if this.Current == nil {
		return nil, false
	}

	next, ok := this.Current.Next(code)
	if ok {
		this.Current = next
	}

	return next, ok
}
