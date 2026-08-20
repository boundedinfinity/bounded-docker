package state

import (
	"fmt"
	"strings"
)

func New(config MachineConfig) (*Machine, error) {
	m := Machine{
		transitions: []*State{},
		keys:        []*Key{},
		capture:     strings.Builder{},
		selected:    []string{},
	}

	// process predefined keys
	for _, kconfig := range config.Keys {
		key := &Key{
			Code: norm(kconfig.Code),
			Name: norm(kconfig.Name),
		}

		if kconfig.Name == "" {
			key.Name = key.Code
		}

		m.keys = append(m.keys, key)
	}

	// process keys in transisitons and navigations
	for _, sconfig := range config.States {
		process := func(km map[string][]string) {
			for _, keys := range km {
				for _, key := range keys {
					key = norm(key)
					if _, ok := m.FindKey(key); !ok {
						m.keys = append(m.keys, &Key{
							Code: key,
							Name: key,
						})
					}
				}
			}
		}

		process(sconfig.Transitions)
		process(sconfig.Navigations)
	}

	for i, sconfig := range config.States {
		sid := norm(sconfig.Id)
		if _, sok := m.FindState(sid); sok {
			return nil, errFn("state[%d:%s]: duplicate", i, sid)
		}

		state := &State{
			Id:          sid,
			Name:        norm(sconfig.Name),
			Transitions: []Transistion{},
		}

		m.transitions = append(m.transitions, state)
	}

	for s, sconfig := range config.States {
		sid := norm(sconfig.Id)
		state, sok := m.FindState(sid)

		if !sok {
			return nil, errFn("state[%d:%s]: not found", s, sid)
		}

		checkDups := func(keys []*Key) error {
			km := map[string][]int{}

			for k, key := range keys {
				if _, ok := km[key.Code]; !ok {
					km[key.Code] = []int{k}
				} else {
					km[key.Code] = append(km[key.Code], k)
				}
			}

			for k, is := range km {
				if len(is) > 1 {
					var found []string
					for _, i := range is {
						found = append(found, fmt.Sprintf("%d", i))
					}
					list := strings.Join(found, ",")

					return errFn("state[%d:%s].transitions[].key[%s]: duplicated at [%s]", s, sid, k, list)
				}
			}

			return nil
		}

		for tid, ckeys := range sconfig.Transitions {
			tid = norm(tid)
			tstate, tok := m.FindState(tid)

			if !tok {
				return nil, errFn("state[%d:%s].transitions[%s]: not found", s, sid, tid)
			}

			transision := Transistion{State: tstate}

			for k, ckey := range ckeys {
				ckey = norm(ckey)
				key, kok := m.FindKey(ckey)

				if !kok {
					return nil, errFn("state[%d:%s].transitions[%s].key[%d:%s]: not found", s, sid, tid, k, ckey)
				}

				transision.Keys = append(transision.Keys, key)
			}

			if err := checkDups(transision.Keys); err != nil {
				return nil, err
			}

			state.Transitions = append(state.Transitions, transision)
		}

		for name, nkeys := range sconfig.Navigations {
			name = norm(name)
			navigation := Navigation{Name: name, Keys: []*Key{}}

			for k, nkey := range nkeys {
				nkey = norm(nkey)
				key, kok := m.FindKey(nkey)

				if !kok {
					return nil, errFn("state[%d:%s].navigation[%s].key[%d:%s]: not found", s, sid, name, k, nkey)
				}

				navigation.Keys = append(navigation.Keys, key)
			}

			if err := checkDups(navigation.Keys); err != nil {
				return nil, err
			}

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
