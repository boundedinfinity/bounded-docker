package state

import (
	"strings"
)

type StateKind int

type stateKinds struct {
	Navigation StateKind
	Capture    StateKind
	Selection  StateKind
}

var StateKinds = stateKinds{
	Navigation: 1,
	Capture:    2,
	Selection:  3,
}

type State struct {
	Id          string
	Name        string
	Kind        StateKind
	Navigations []Navigation
	Transitions []Transistion
	selected    []string
	capture     strings.Builder
}

func (this *State) Next(code string) (*State, bool) {
	for _, transistion := range this.Transitions {
		if transistion.Matches(code) {
			return transistion.State, true
		}
	}

	for _, navigation := range this.Navigations {
		if navigation.Matches(code) {
			return this, true
		}
	}

	return this, false
}

// func (this State) HelpView() []key.Binding {
// 	bindings := make([]key.Binding, 0)

// 	bindings = append(bindings, this.TransisionHelp()...)
// 	bindings = append(bindings, this.NavigationHelp()...)

// 	return bindings
// }

// func (this State) HelpView2() [][]key.Binding {
// 	bindings := [][]key.Binding{
// 		this.TransisionHelp(),
// 		this.NavigationHelp(),
// 	}

// 	return bindings
// }

// func (this State) TransisionHelp() []key.Binding {
// 	bindings := make([]key.Binding, len(this.Transitions))

// 	for i, transition := range this.Transitions {
// 		bindings[i] = transition.bindings
// 	}

// 	return bindings
// }

// func (this State) NavigationHelp() []key.Binding {
// 	bindings := make([]key.Binding, len(this.Navigations))

// 	for i, navigation := range this.Navigations {
// 		bindings[i] = navigation.bindings
// 	}

// 	return bindings
// }
