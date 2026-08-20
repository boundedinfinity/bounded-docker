package state

type navigationKind int

const (
	_ navigationKind = iota
	transision
	navigation
)

type Navigation struct {
	Name string
	Keys []*Key
	kind navigationKind
}

func (this Navigation) Matches(code string) bool {
	for _, key := range this.Keys {
		if key.Matches(code) {
			return true
		}
	}

	return false
}
