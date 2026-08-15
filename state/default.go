package state

func DefaultConfig() []ConfigState {
	return []ConfigState{
		{
			Name:        "root",
			Description: "",
			Bindings: []ConfigBinding{
				{Keys: []ConfigKey{{Code: "c"}}, Next: "containers"},
				{Keys: []ConfigKey{{Code: "i"}}, Next: "images"},
				{Keys: []ConfigKey{{Code: "e"}}, Next: "errors"},
			},
		},
		{
			Name:        "containers",
			Description: "List of containers",
			Bindings: []ConfigBinding{
				Back("root"),
				Up(),
				Down(),
			},
		},
		{
			Name:        "images",
			Description: "List of images",
			Bindings: []ConfigBinding{
				Back("root"),
				Up(),
				Down(),
			},
		},
		{
			Name:        "errors",
			Description: "List of errors",
			Bindings: []ConfigBinding{
				Back("root"),
				Up(),
				Down(),
			},
		},
	}
}
