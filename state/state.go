package state

type State struct {
	Id          string
	Name        string
	Commands    []Command
	Transitions []Transistion
}

func (this *State) Next(code string) (*State, *Command, bool) {
	for _, transistion := range this.Transitions {
		if transistion.Matches(code) {
			return transistion.State, nil, true
		}
	}

	for _, command := range this.Commands {
		if command.Matches(code) {
			return this, &command, true
		}
	}

	return this, nil, false
}
