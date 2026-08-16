package state

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	ErrMachine = errorer.New("state machine")
	ErrConfig  = errorer.New("config")
	errFn      = errorer.Func(ErrMachine, ErrConfig)
)

// /////////////////////////////////////////////////////////////////////////////////////////////////

type Machine struct {
	Current *State
	states  []*State
	keys    []*Key
}

func (this Machine) FindState(id string) (*State, bool) {
	for _, state := range this.states {
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

func (this *Machine) Next(k tea.KeyPressMsg) (*State, bool) {
	if this.Current == nil {
		return nil, false
	}

	next, ok := this.Current.Next(k)
	if ok {
		this.Current = next
	}

	return next, ok
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type State struct {
	Id          string
	Name        string
	Navigations []Navigation
	Transitions []Transistion
}

func (this *State) Next(k tea.KeyPressMsg) (*State, bool) {
	for _, transistion := range this.Transitions {
		if transistion.Matches(k) {
			return transistion.State, true
		}
	}

	for _, navigation := range this.Navigations {
		if navigation.Matches(k) {
			return this, true
		}
	}

	return this, false
}

func (this State) HelpView() []key.Binding {
	bindings := make([]key.Binding, 0)

	for _, transition := range this.Transitions {
		bindings = append(bindings, transition.bindings)
	}

	for _, navigation := range this.Navigations {
		bindings = append(bindings, navigation.bindings)
	}

	return bindings
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type Key struct {
	Code string
	Name string
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type Transistion struct {
	State    *State
	Keys     []*Key
	bindings key.Binding
}

func (this Transistion) Matches(k tea.KeyPressMsg) bool {
	return key.Matches(k, this.bindings)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type Navigation struct {
	Name     string
	Keys     []*Key
	bindings key.Binding
}

func (this Navigation) Matches(k tea.KeyPressMsg) bool {
	return key.Matches(k, this.bindings)
}
