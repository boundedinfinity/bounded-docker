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
	Current   *State
	states    map[string]*State
	commands  map[string]*Command
	keys      map[string]*Key
	startFunc []func()
	endFunc   []func()
}

func (this *Machine) GetState(id string) (*State, bool) {
	state, ok := this.states[id]
	return state, ok
}

func (this *Machine) GetCommand(id string) (*Command, bool) {
	command, ok := this.commands[id]
	return command, ok
}

func (this Machine) States() []*State {
	return mapper.Values(this.states)
}

func (this Machine) StartState() *State {
	return this.Current
}

func (this *Machine) Next(code string) (*State, []*Command) {
	if this.Current == nil {
		panic("state machine is not initialized")
	}

	var state *State
	var commands []*Command

	state, commands = this.Current.Next(code)

	if state != nil && state.Id != this.Current.Id {
		this.Current.Exit()
		this.Current = state
		this.Current.Enter()
	}

	for _, c := range commands {
		if c != nil {
			c.Run(this.Current)
		}
	}

	return this.Current, commands
}

func (this *Machine) AddStartFunc(fn func()) {
	if fn != nil {
		this.startFunc = append(this.startFunc, fn)
	}
}

func (this *Machine) AddEndFunc(fn func()) {
	if fn != nil {
		this.endFunc = append(this.endFunc, fn)
	}
}

func (this Machine) Start() {
	for _, fn := range this.startFunc {
		if fn != nil {
			fn()
		}
	}
}

func (this Machine) End() {
	for _, fn := range this.endFunc {
		if fn != nil {
			fn()
		}
	}
}
