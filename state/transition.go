package state

type Transistion struct {
	State *State
	Keys  []*Key
}

func (this Transistion) Matches(code string) bool {
	for _, key := range this.Keys {
		if key.Matches(code) {
			return true
		}
	}
	return false
}
