package state

type ConfigState struct {
	Name        string
	Bindings    []ConfigBinding
	Description string
}

type ConfigBinding struct {
	Keys []ConfigKey
	Help string
	Next string
}

type ConfigKey struct {
	Code string
	Help string
}

func Back(name string) ConfigBinding {
	return ConfigBinding{
		Help: "go back",
		Next: name,
		Keys: []ConfigKey{
			{Code: "left", Help: "←"},
			{Code: "h"},
		},
	}
}

func Up() ConfigBinding {
	return ConfigBinding{
		Help: "move up",
		Keys: []ConfigKey{
			{Code: "up", Help: "↑"},
			{Code: "k"},
		},
	}
}

func Down() ConfigBinding {
	return ConfigBinding{
		Help: "move down",
		Keys: []ConfigKey{
			{Code: "down", Help: "↓"},
			{Code: "j"},
		},
	}
}
