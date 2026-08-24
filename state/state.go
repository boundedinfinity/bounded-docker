package state

type State struct {
	Id          string
	Name        string
	Description string
	Commands    []StateCommand
	Transitions []Transistion
	enterFunc   []func(*State)
	exitFunc    []func(*State)
}

func (this *State) Next(code string) (*State, []*Command) {
	var state *State
	var commands []*Command

	for _, transistion := range this.Transitions {
		if transistion.Matches(code) {
			state = transistion.State
			break
		}
	}

	for _, stateCommand := range this.Commands {
		if stateCommand.Matches(code) {
			commands = append(commands, stateCommand.Command)
		}
	}

	return state, commands
}

func (this *State) AddEnterFunc(fn func(*State)) {
	if fn != nil {
		this.enterFunc = append(this.enterFunc, fn)
	}
}

func (this *State) Enter() {
	for _, fn := range this.enterFunc {
		if fn != nil {
			fn(this)
		}
	}
}

func (this *State) AddExitFunc(fn func(*State)) {
	if fn != nil {
		this.exitFunc = append(this.exitFunc, fn)
	}
}

func (this *State) Exit() {
	for _, fn := range this.exitFunc {
		if fn != nil {
			fn(this)
		}
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type Transistion struct {
	State *State
	Keys  KeyList
}

func (this Transistion) Matches(code string) bool {
	return this.Keys.Matches(code)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type Command struct {
	Id          string
	Name        string
	Description string
	keys        KeyList
	runFunc     []func(*State, *Command)
}

func (this Command) Matches(code string) bool {
	return this.keys.Matches(code)
}

func (this *Command) AddRunFunc(fn func(*State, *Command)) {
	if fn != nil {
		this.runFunc = append(this.runFunc, fn)
	}
}

func (this Command) Run(s *State) {
	for _, fn := range this.runFunc {
		if fn != nil {
			fn(s, &this)
		}
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type StateCommand struct {
	includeDefaults bool
	Command         *Command
	keys            KeyList
}

func (this StateCommand) Matches(code string) bool {
	if this.includeDefaults && len(this.Command.keys) > 0 {
		if this.Command.Matches(code) {
			return true
		}
	}

	return this.keys.Matches(code)
}

func (this StateCommand) Keys() KeyList {
	var keys KeyList
	if this.includeDefaults && len(this.Command.keys) > 0 {
		keys = append(keys, this.Command.keys...)
	}

	keys = append(keys, this.keys...)
	return keys
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type KeyList []*Key

func (this KeyList) Matches(code string) bool {
	for _, key := range this {
		if key.Matches(code) {
			return true
		}
	}

	return false
}

/////////////////////////////////////////////////////////////////////////////////////////////////

type Key struct {
	Code string
	Name string
}

func (this Key) Matches(code string) bool {
	return this.Code == code
}

var (
	validKeyCodes = []string{
		"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"up", "down", "left", "right",
		"esc", "enter", "space", "tab", "backspace", "delete", "insert",
		"home", "end", "pageup", "pagedown",
	}
)
