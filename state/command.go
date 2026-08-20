package state

type Command struct {
	Name string
	Keys []*Key
}

func (this Command) Matches(code string) bool {
	for _, key := range this.Keys {
		if key.Matches(code) {
			return true
		}
	}

	return false
}
