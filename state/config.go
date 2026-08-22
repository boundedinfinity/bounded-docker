package state

type MachineConfig struct {
	Start       string                   `json:"start"`
	States      map[string]StateConfig   `json:"states"`
	Transitions map[string][]string      `json:"transitions"`
	Commands    map[string]CommandConfig `json:"commands"`
	Keys        map[string]KeyConfig     `json:"keys"`
}

type StateConfig struct {
	Name        string              `json:"name"`
	Commands    map[string][]string `json:"commands"`
	Keys        map[string][]string `json:"keys"`
	Transitions map[string][]string `json:"transitions"`
}

type CommandConfig struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keys        []string `json:"keys"`
}

type KeyConfig struct {
	Name string `json:"name"`
}
