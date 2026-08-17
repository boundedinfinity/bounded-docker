package state

import tea "charm.land/bubbletea/v2"

// /////////////////////////////////////////////////////////////////////////////////////////////////

type StateChangeMsg struct {
	StateId string
}

func StateChangeCmd(id string) func() tea.Msg {
	return func() tea.Msg {
		return StateChangeMsg{
			StateId: id,
		}
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type StateChangedCapture struct {
	StateId string
	Capture string
}

func StateChangedCaptureCmd(id, capture string) func() tea.Msg {
	return func() tea.Msg {
		return StateChangedCapture{
			StateId: id,
			Capture: capture,
		}
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type StateChangedSelected struct {
	StateId  string
	Selected []string
}

func StateChangedSelectedCmd(id string, selected ...string) func() tea.Msg {
	return func() tea.Msg {
		return StateChangedSelected{
			StateId:  id,
			Selected: selected,
		}
	}
}
