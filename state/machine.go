package state

import (
	"strings"

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
	Current  *State
	states   []*State
	keys     []*Key
	capture  strings.Builder
	selected []string
	models   map[string]tea.Model
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

func (this *Machine) Broadcast(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	for _, model := range this.models {
		_, cmd := model.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (this *Machine) Update(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if next, nok := this.Next(msg); nok {
		if model, vok := this.models[next.Id]; vok {
			return model.Update(msg)
		}
	}

	return this.models[this.Current.Id].Update(msg)
}

func (this *Machine) Next(msg tea.KeyPressMsg) (*State, bool) {
	if this.Current == nil {
		return nil, false
	}

	next, ok := this.Current.Next(msg)
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

	bindings = append(bindings, this.TransisionHelp()...)
	bindings = append(bindings, this.NavigationHelp()...)

	return bindings
}

func (this State) HelpView2() [][]key.Binding {
	bindings := [][]key.Binding{
		this.TransisionHelp(),
		this.NavigationHelp(),
	}

	return bindings
}

func (this State) TransisionHelp() []key.Binding {
	bindings := make([]key.Binding, len(this.Transitions))

	for i, transition := range this.Transitions {
		bindings[i] = transition.bindings
	}

	return bindings
}

func (this State) NavigationHelp() []key.Binding {
	bindings := make([]key.Binding, len(this.Navigations))

	for i, navigation := range this.Navigations {
		bindings[i] = navigation.bindings
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
