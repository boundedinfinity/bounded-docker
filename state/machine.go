package state

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	ErrStateMachine = errorer.New("state machine")
	ErrConfig       = errorer.New("config")
	errFn           = errorer.Func(ErrStateMachine, ErrConfig)
)

type Machine struct {
	m       map[string]*State
	Current *State
}

func (this *Machine) Next(k tea.KeyPressMsg) (*State, bool) {
	next, nok := this.Current.Next(k)

	if nok {
		if state, sok := this.m[next]; sok {
			this.Current = state
			return this.Current, sok
		}
	}

	return this.Current, nok
}

type State struct {
	Name     string
	bindings []stateBinding
	Help     string
}

func (this State) Next(k tea.KeyPressMsg) (string, bool) {
	for _, sbinding := range this.bindings {
		b := sbinding.binding
		if key.Matches(k, b) {
			return sbinding.next, true
		}
	}

	return "", false
}

func (this State) HelpView() []key.Binding {
	bindings := make([]key.Binding, len(this.bindings))

	for i := range this.bindings {
		bindings[i] = this.bindings[i].binding
	}

	return bindings
}

type stateBinding struct {
	binding key.Binding
	next    string
}
