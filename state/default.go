package state

func DefaultConfig() MachineConfig {
	return MachineConfig{
		Start: "root",
		States: map[string]StateConfig{
			"root": {
				Name: "Dashboard",
				Transitions: map[string][]string{
					"container.list": {"c"},
					"image.list":     {"i"},
					"network.list":   {"n"},
					"errors":         {"e"},
				},
				Commands: map[string][]string{
					"quit": {"[command]", "esc"},
				},
			},
			"errors": {
				Name: "Errors",
				Transitions: map[string][]string{
					"root":   {"r", "left", "h", "b", "esc"},
					"errors": {"y"},
				},
				Commands: map[string][]string{
					"quit":      {},
					"line.up":   {},
					"line.down": {},
				},
			},
			"container.list": {
				Name: "Containers",
				Commands: map[string][]string{
					"quit":               {},
					"line.up":            {"up", "j"},
					"line.down":          {"down", "k"},
					"container.list.all": {"a"},
				},
				Transitions: map[string][]string{
					"root":              {"r", "left", "h", "b", "esc"},
					"container.details": {"d", "enter"},
					"container.logs":    {"l"},
				},
			},
			"container.details": {
				Name: "Container Details",
				Transitions: map[string][]string{
					"root":           {"r"},
					"container.list": {"c", "left", "h", "b", "esc"},
				},
				Commands: map[string][]string{
					"quit": {"[command]"},
				},
			},
			"container.logs": {
				Name: "Container Logs",
				Transitions: map[string][]string{
					"root":           {"r"},
					"container.list": {"c", "left", "h", "b", "esc"},
				},
				Commands: map[string][]string{
					"container.logs.follow": {"f"},
					"quit":                  {"[command]"},
				},
			},
			"image.list": {
				Name: "Images",
				Transitions: map[string][]string{
					"root":          {"r", "left", "h", "b", "esc"},
					"image.details": {"d", "enter"},
				},
				Commands: map[string][]string{
					"up":             {"up", "j"},
					"down":           {"down", "k"},
					"image.list.all": {"a"},
					"quit":           {},
				},
			},
			"image.details": {
				Name: "Image Details",
				Transitions: map[string][]string{
					"root":       {"r"},
					"image.list": {"i", "left", "h", "b", "esc"},
				},
				Commands: map[string][]string{
					"quit": {},
				},
			},
			"network.list": {
				Name: "Networks",
				Transitions: map[string][]string{
					"root":            {"r", "left", "h", "b", "esc"},
					"network.details": {"d", "enter"},
				},
				Commands: map[string][]string{
					"up":   {"up", "j"},
					"down": {"down", "k"},
					"quit": {},
				},
			},
			"network.details": {
				Name: "Network Details",
				Transitions: map[string][]string{
					"root":         {"r"},
					"network.list": {"n", "left", "h", "b", "esc"},
				},
				Commands: map[string][]string{
					"quit": {},
				},
			},
		},
		Commands: map[string]CommandConfig{
			"quit": {
				Name: "Quit",
				Keys: []string{"q"},
			},
			"line.up": {
				Name: "Move Up",
				Keys: []string{"up", "j"},
			},
			"line.down": {
				Name: "Move Down",
				Keys: []string{"down", "k"},
			},
			"container.list.all": {
				Name: "List All Containers",
			},
			"container.logs.follow": {
				Name: "Follow Container Logs",
			},
			"image.list.all": {
				Name: "List All Images",
			},
			"errors.copy": {
				Name: "Copy error to clipboard",
				Keys: []string{"c"},
			},
		},
		Keys: map[string]KeyConfig{
			"left":  {Name: "←"},
			"up":    {Name: "↑"},
			"down":  {Name: "↓"},
			"right": {Name: "→"},
			"enter": {Name: "↵"},
		},
	}
}
