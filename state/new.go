package state

import "errors"

func New(config MachineConfig) (*Machine, error) {
	m := Machine{
		states:   map[string]*State{},
		commands: map[string]*Command{},
		keys:     map[string]*Key{},
	}

	getKey := func(code string) (*Key, error) {
		if !Utils.K.valid(code) {
			return nil, errors.New("invalid key")
		}

		if _, ok := m.keys[code]; !ok {
			m.keys[code] = &Key{Code: code, Name: code}
		}

		return m.keys[code], nil
	}

	getCommand := func(id string) *Command {
		if _, ok := m.commands[id]; !ok {
			m.commands[id] = &Command{Id: id, Name: id, Description: id}
		}

		return m.commands[id]
	}

	for code, config := range config.Keys {
		key := &Key{
			Code: Utils.S.norm(code),
			Name: Utils.S.norm(config.Name),
		}

		if _, ok := m.keys[key.Code]; ok {
			return nil, errFn("keys[%s]: duplicate", key.Code)
		}

		m.keys[key.Code] = key
	}

	for id, config := range config.Commands {
		command := &Command{
			Id:          Utils.S.norm(id),
			Name:        Utils.S.norm(config.Name),
			Description: Utils.S.norm(config.Description),
			keys:        KeyList{},
		}

		if _, ok := m.commands[command.Id]; ok {
			return nil, errFn("command[%s]: duplicate", command.Id)
		}

		for i, keyCode := range config.Keys {
			keyCode = Utils.S.norm(keyCode)

			switch keyCode {
			case "[command]", "[defaults]":
				return nil, errFn("command[%s][%d]: invalid key code '%s'", command.Id, i, keyCode)
			default:
				if !Utils.K.valid(keyCode) {
					return nil, errFn("command[%s][%d]: invalid key code '%s'", command.Id, i, keyCode)
				}
			}

			if key, err := getKey(keyCode); err == nil {
				command.keys = append(command.keys, key)
			} else {
				return nil, errFn("commands[%s][%d]: %v", command.Id, i, err)
			}
		}

		m.commands[command.Id] = command
	}

	for id, config := range config.States {
		state := &State{
			Id:          Utils.S.norm(id),
			Name:        Utils.S.norm(config.Name),
			Transitions: []Transistion{},
			Commands:    []StateCommand{},
		}

		if _, ok := m.states[state.Id]; ok {
			return nil, errFn("state[%s]: duplicate", state.Id)
		}

		m.states[state.Id] = state
	}

	for stateId, config := range config.States {
		state, ok := m.states[stateId]
		if !ok {
			return nil, errFn("state[%s]: not found", stateId)
		}

		for trantsionId, keyCodes := range config.Transitions {
			transistionState, ok := m.states[trantsionId]
			if !ok {
				return nil, errFn("state[%s].transitions[%s]: not found", stateId, trantsionId)
			}

			transition := Transistion{State: transistionState, Keys: KeyList{}}

			for i, keyCode := range keyCodes {
				keyCode = Utils.S.norm(keyCode)
				if key, err := getKey(keyCode); err == nil {
					transition.Keys = append(transition.Keys, key)
				} else {
					return nil, errFn("state[%s].transitions[%s][%d]: %v", stateId, trantsionId, i, err)
				}
			}

			state.Transitions = append(state.Transitions, transition)
		}

		for commandId, keyCodes := range config.Commands {
			command := getCommand(commandId)
			stateCommand := StateCommand{Command: command, keys: KeyList{}}

			if len(keyCodes) == 0 {
				stateCommand.includeDefaults = true
			} else {
				for i, keyCode := range keyCodes {
					keyCode = Utils.S.norm(keyCode)

					switch keyCode {
					case "[command]", "[defaults]":
						stateCommand.includeDefaults = true
					default:
						if key, err := getKey(keyCode); err == nil {
							stateCommand.keys = append(stateCommand.keys, key)
						} else {
							return nil, errFn("state[%s].commands[%s][%d]: %v", stateId, commandId, i, err)
						}
					}
				}
			}

			state.Commands = append(state.Commands, stateCommand)
		}
	}

	for stateId, state := range m.states {
		validator := stateKeyValidator{}

		for _, transition := range state.Transitions {
			for k, key := range transition.Keys {
				validator.add(stateId, transition.State.Id, "", key.Code, k)
			}
		}

		for _, command := range state.Commands {
			for k, key := range command.Keys() {
				validator.add(stateId, "", command.Command.Id, key.Code, k)
			}
		}

		if err := validator.validate(); err != nil {
			return nil, err
		}
	}

	if start, ok := m.states[config.Start]; !ok {
		return nil, errFn("start state '%s' does not exist", config.Start)
	} else {
		m.Current = start
	}

	return &m, nil
}
