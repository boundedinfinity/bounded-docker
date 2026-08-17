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
	Current     *State
	transitions []*State
	navigation  []*State
	keys        []*Key
	capture     strings.Builder
	selected    []string
	models      map[string]tea.Model
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

type navigationKind int

const (
	_ navigationKind = iota
	transision
	navigation
)

type Navigation struct {
	Name     string
	Keys     []*Key
	kind     navigationKind
	bindings key.Binding
}

func (this Navigation) Matches(k tea.KeyPressMsg) bool {
	return key.Matches(k, this.bindings)
}
