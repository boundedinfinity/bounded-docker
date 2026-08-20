package state

type State struct {
	Id          string
	Name        string
	Commands    []Command
	Transitions []Transistion
}

func (this *State) Next(code string) (*State, bool) {
	for _, transistion := range this.Transitions {
		if transistion.Matches(code) {
			return transistion.State, true
		}
	}

	for _, navigation := range this.Commands {
		if navigation.Matches(code) {
			return this, true
		}
	}

	return this, false
}
