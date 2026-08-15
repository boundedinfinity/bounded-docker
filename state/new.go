package state

import (
	"strings"

	"charm.land/bubbles/v2/key"
)

func NewMachine(start string, stateConfigs []ConfigState) (*Machine, error) {
	machine := Machine{
		m: make(map[string]*State),
	}

	for s, stateConfig := range stateConfigs {
		state := State{
			Name:     stateConfig.Name,
			Help:     stateConfig.Description,
			bindings: []stateBinding{},
		}

		for b, bindingConfig := range stateConfig.Bindings {
			if len(bindingConfig.Keys) < 1 {
				return nil, errFn("state[ %d].binding[%d]: no keys defined", s, b)
			}

			codes := make([]string, len(bindingConfig.Keys))
			helps := make([]string, len(bindingConfig.Keys))

			for k, kconfig := range bindingConfig.Keys {
				code := strings.TrimSpace(kconfig.Code)
				help := strings.TrimSpace(kconfig.Help)
				if help == "" {
					help = code
				}

				codes[k] = code
				helps[k] = help
			}

			help := strings.Join(helps, "/")

			sbinding := stateBinding{
				binding: key.NewBinding(
					key.WithKeys(codes...),
					key.WithHelp(help, bindingConfig.Help),
				),
				next: bindingConfig.Next,
			}

			state.bindings = append(state.bindings, sbinding)
		}

		machine.m[state.Name] = &state
	}

	for _, state := range machine.m {
		for _, binding := range state.bindings {
			if binding.next == "" {
				continue
			}

			if _, ok := machine.m[binding.next]; !ok {
				return nil, errFn("state[%s].binding[%s].next '%s' does not exist", state.Name, binding.binding.Help().Key, binding.next)
			}
		}
	}

	if _, ok := machine.m[start]; !ok {
		return nil, errFn("start state '%s' does not exist", start)
	}

	machine.Current = machine.m[start]
	return &machine, nil
}
