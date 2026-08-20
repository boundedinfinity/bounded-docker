package state

func DefaultConfig() MachineConfig {
	return MachineConfig{
		Start: "root",
		States: []StateConfig{
			{
				Id:   "quit",
				Name: "Quit",
			},
			{
				Id:   "root",
				Name: "Dashboard",
				Transitions: map[string][]string{
					"containers": {"c"},
					"images":     {"i"},
					"networks":   {"n"},
					"errors":     {"e"},
					"quit":       {"q", "esc"},
				},
			},
			{
				Id:   "containers",
				Name: "Containers",
				Commands: map[string][]string{
					"up":   {"up", "j"},
					"down": {"down", "k"},
				},
				Transitions: map[string][]string{
					"root":              {"r", "left", "h", "b", "esc"},
					"container-details": {"d", "enter"},
					"quit":              {"q"},
				},
			},
			{
				Id:   "container-details",
				Name: "Container Details",
				Transitions: map[string][]string{
					"root":       {"r"},
					"containers": {"c", "left", "h", "b", "esc"},
					"quit":       {"q"},
				},
			},
			{
				Id:   "images",
				Name: "Images",
				Transitions: map[string][]string{
					"root":          {"r", "left", "h", "b", "esc"},
					"image-details": {"d", "enter"},
					"quit":          {"q"},
				},
				Commands: map[string][]string{
					"up":   {"up", "j"},
					"down": {"down", "k"},
				},
			},
			{
				Id:   "image-details",
				Name: "Image Details",
				Transitions: map[string][]string{
					"root":   {"r"},
					"images": {"i", "left", "h", "b", "esc"},
					"quit":   {"q"},
				},
			},
			{
				Id:   "networks",
				Name: "Networks",
				Transitions: map[string][]string{
					"root":            {"r", "left", "h", "b", "esc"},
					"network-details": {"d", "enter"},
					"quit":            {"q"},
				},
				Commands: map[string][]string{
					"up":   {"up", "j"},
					"down": {"down", "k"},
				},
			},
			{
				Id:   "network-details",
				Name: "Network Details",
				Transitions: map[string][]string{
					"root":     {"r"},
					"networks": {"n", "left", "h", "b", "esc"},
					"quit":     {"q"},
				},
			},
			{
				Id:   "errors",
				Name: "Errors",
				Transitions: map[string][]string{
					"root":   {"r", "left", "h", "b", "esc"},
					"quit":   {"q"},
					"errors": {"y"},
				},
				Commands: map[string][]string{
					"up":   {"up", "j"},
					"down": {"down", "k"},
				},
			},
		},
		Keys: []KeyConfig{
			{Code: "left", Name: "←"},
			{Code: "up", Name: "↑"},
			{Code: "down", Name: "↓"},
			{Code: "right", Name: "→"},
			{Code: "enter", Name: "↵"},
		},
	}
}
