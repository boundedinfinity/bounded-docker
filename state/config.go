package state

type MachineConfig struct {
	Start  string        `json:"start"`
	States []StateConfig `json:"states"`
	Keys   []KeyConfig   `json:"keys"`
}

type StateConfig struct {
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	Navigations map[string][]string `json:"navigations"`
	Transitions map[string][]string `json:"transitions"`
}

type TransistionConfig interface {
	map[string][]string | []string
}

type KeyConfig struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
