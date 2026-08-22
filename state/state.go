package state

type State struct {
	Id          string
	Name        string
	Description string
	Commands    []StateCommand
	Transitions []Transistion
}

func (this *State) Next(code string) (*State, *Command) {
	for _, transistion := range this.Transitions {
		if transistion.Matches(code) {
			return transistion.State, nil
		}
	}

	for _, command := range this.Commands {
		if command.Matches(code) {
			return this, command.Command
		}
	}

	return this, nil
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
}

func (this Command) Matches(code string) bool {
	return this.keys.Matches(code)
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
