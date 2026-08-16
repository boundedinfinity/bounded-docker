package state

import (
	"strings"

	bkey "charm.land/bubbles/v2/key"
)

func New(config MachineConfig) (*Machine, error) {
	n := func(s string) string {
		return strings.TrimSpace(s)
	}

	m := Machine{
		states: []*State{},
		keys:   []*Key{},
	}

	for _, kconfig := range config.Keys {
		key := &Key{
			Code: n(kconfig.Code),
			Name: n(kconfig.Name),
		}

		if kconfig.Name == "" {
			key.Name = key.Code
		}

		m.keys = append(m.keys, key)
	}

	for _, sconfig := range config.States {
		for _, keys := range sconfig.Transitions {
			for _, key := range keys {
				key = n(key)
				if _, ok := m.FindKey(key); !ok {
					m.keys = append(m.keys, &Key{
						Code: key,
						Name: key,
					})
				}
			}
		}

		for _, keys := range sconfig.Navigations {
			for _, key := range keys {
				key = n(key)
				if _, ok := m.FindKey(key); !ok {
					m.keys = append(m.keys, &Key{
						Code: key,
						Name: key,
					})
				}
			}
		}
	}

	for i, sconfig := range config.States {
		sid := n(sconfig.Id)
		if _, sok := m.FindState(sid); sok {
			return nil, errFn("state[%d:%s]: duplicate", i, sid)
		}

		state := &State{
			Id:          sid,
			Name:        n(sconfig.Name),
			Transitions: []Transistion{},
		}

		m.states = append(m.states, state)
	}

	for s, sconfig := range config.States {
		sid := n(sconfig.Id)
		state, sok := m.FindState(sid)

		if !sok {
			return nil, errFn("state[%d:%s]: not found", s, sid)
		}

		for tid, ckeys := range sconfig.Transitions {
			tid = n(tid)
			tstate, tok := m.FindState(tid)

			if !tok {
				return nil, errFn("state[%d:%s].transition[%s]: not found", s, sid, tid)
			}

			transision := Transistion{
				State: tstate,
			}

			for k, ckey := range ckeys {
				ckey = n(ckey)
				key, kok := m.FindKey(ckey)

				if !kok {
					return nil, errFn("state[%d:%s].transition[%s].key[%d:%s]: not found", s, sid, tid, k, ckey)
				}

				transision.Keys = append(transision.Keys, key)
			}

			state.Transitions = append(state.Transitions, transision)
		}

		for name, nkeys := range sconfig.Navigations {
			name = n(name)
			navigation := Navigation{
				Name: name,
				Keys: []*Key{},
			}

			for k, nkey := range nkeys {
				nkey = n(nkey)
				key, kok := m.FindKey(nkey)

				if !kok {
					return nil, errFn("state[%d:%s].navigation[%s].key[%d:%s]: not found", s, sid, name, k, nkey)
				}

				navigation.Keys = append(navigation.Keys, key)
			}

			state.Navigations = append(state.Navigations, navigation)
		}
	}

	for _, state := range m.states {
		for _, transistion := range state.Transitions {
			names := make([]string, len(transistion.Keys))
			codes := make([]string, len(transistion.Keys))

			for k, key := range transistion.Keys {
				names[k] = key.Name
				codes[k] = key.Code
			}

			help := strings.Join(names, "/")

			transistion.bindings = bkey.NewBinding(
				bkey.WithKeys(codes...),
				bkey.WithHelp(help, transistion.State.Name),
			)

			state.Transitions = append(state.Transitions, transistion)
		}

		for _, navigation := range state.Navigations {
			names := make([]string, len(navigation.Keys))
			codes := make([]string, len(navigation.Keys))

			for k, key := range navigation.Keys {
				names[k] = key.Name
				codes[k] = key.Code
			}

			help := strings.Join(names, "/")

			navigation.bindings = bkey.NewBinding(
				bkey.WithKeys(codes...),
				bkey.WithHelp(help, navigation.Name),
			)

			state.Navigations = append(state.Navigations, navigation)
		}
	}

	if start, ok := m.FindState(config.Start); !ok {
		return nil, errFn("start state '%s' does not exist", config.Start)
	} else {
		m.Current = start
	}

	return &m, nil
}
