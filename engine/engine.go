package engine

import "slices"

type State interface {
	Name() string
	Help() string
	Handles(event string) bool
	Run(event string) (State, Command)
}

func matchState(event string) func(State) bool {
	return func(state State) bool {
		return state.Handles(event)
	}
}

type Command interface {
}

type Engine struct {
}

// ////////////////////////////////////////////////////////////////////////////////////////////////
// CaptureState
// ////////////////////////////////////////////////////////////////////////////////////////////////

var _ State = &CaptureState{}

type CaptureState struct {
	name   string
	help   string
	key    string
	text   string
	parent State
}

// Handles implements [State].
func (c *CaptureState) Handles(event string) bool {
	panic("unimplemented")
}

// Help implements [State].
func (c *CaptureState) Help() string {
	panic("unimplemented")
}

// Name implements [State].
func (c *CaptureState) Name() string {
	panic("unimplemented")
}

// Run implements [State].
func (c *CaptureState) Run(event string) (State, Command) {
	panic("unimplemented")
}

// ////////////////////////////////////////////////////////////////////////////////////////////////
// NavigationState
// ////////////////////////////////////////////////////////////////////////////////////////////////

var _ State = &NavigationState{}

type NavigationState struct {
	name   string
	help   string
	key    string
	next   []State
	parent State
}

// Handles implements [State].
func (this *NavigationState) Handles(event string) bool {
	panic("unimplemented")
}

// Help implements [State].
func (this *NavigationState) Help() string {
	panic("unimplemented")
}

// Name implements [State].
func (this *NavigationState) Name() string {
	panic("unimplemented")
}

// Run implements [State].
func (this *NavigationState) Run(event string) (State, Command) {
	panic("unimplemented")
}

// ////////////////////////////////////////////////////////////////////////////////////////////////
// TransistionState
// ////////////////////////////////////////////////////////////////////////////////////////////////

var _ State = &TransistionState{}

type TransistionState struct {
	name   string
	help   string
	key    string
	next   []State
	parent State
}

func (this *TransistionState) Help() string {
	panic("unimplemented")
}

func (this *TransistionState) Name() string {
	panic("unimplemented")
}

func (this *TransistionState) Handles(event string) bool {
	return slices.ContainsFunc(this.next, matchState(event))
}

func (this *TransistionState) Run(event string) (State, Command) {
	index := slices.IndexFunc(this.next, matchState(event))
	return this.next[index], nil
}
