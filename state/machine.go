package state

import (
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/idiomatic/mapper"
)

var (
	ErrMachine = errorer.New("state machine")
	ErrConfig  = errorer.New("config")
	errFn      = errorer.Func(ErrMachine, ErrConfig)
)

// /////////////////////////////////////////////////////////////////////////////////////////////////

type Machine struct {
	Current  *State
	states   map[string]*State
	commands map[string]*Command
	keys     map[string]*Key
}

func (this Machine) States() []*State {
	return mapper.Values(this.states)
}

func (this Machine) Start() *State {
	return this.Current
}

func (this *Machine) Next(code string) (*State, *Command) {
	if this.Current == nil {
		panic("state machine is not initialized")
	}

	var cmd *Command
	this.Current, cmd = this.Current.Next(code)

	return this.Current, cmd
}
