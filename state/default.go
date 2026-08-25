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
					"container.inspect": {"i"},
					"container.logs":    {"l", "enter"},
				},
			},
			"container.inspect": {
				Name: "Inspect",
				Transitions: map[string][]string{
					"root":           {"r"},
					"container.list": {"left", "h", "b", "esc"},
					"container.logs": {"l"},
				},
				Commands: map[string][]string{
					"container.inspect.size": {"s"},
					"container.inspect.copy": {"c"},
					"line.up":                {"up", "j"},
					"line.down":              {"down", "k"},
					"quit":                   {"[command]"},
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
					"quit":           {},
					"line.up":        {"up", "j"},
					"line.down":      {"down", "k"},
					"image.list.all": {"a"},
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
				Name: "Toggle List All Containers",
			},
			"container.logs.follow": {
				Name: "Follow Container Logs",
			},
			"container.inspect.size": {
				Name: "Toggle Filesystem Size",
			},
			"container.inspect.copy": {
				Name: "Copy Value to Clipboard",
			},
			"image.list.all": {
				Name: "Toggle List All Images",
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
