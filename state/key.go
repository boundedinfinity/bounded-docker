package state

type Key struct {
	Code string
	Name string
}

func (this Key) Matches(code string) bool {
	return this.Code == code
}
