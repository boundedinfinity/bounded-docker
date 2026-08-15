package state_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/boundedinfinity/docker-tui/state"
	"github.com/stretchr/testify/assert"
)

// v2.KeyPressMsg {Text: "c", Mod: 0, Code: 99, ShiftedCode: 0, BaseCode: 0, IsRepeat: false}
// v2.KeyPressMsg {Text: "e", Mod: 0, Code: 101, ShiftedCode: 0, BaseCode: 0, IsRepeat: false}
// v2.KeyPressMsg {Text: "i", Mod: 0, Code: 105, ShiftedCode: 0, BaseCode: 0, IsRepeat: false}

func Test_StateMachine_Creation(t *testing.T) {
	machine, err := state.NewMachine("root", state.DefaultConfig())
	assert.NoError(t, err)
	assert.Equal(t, "root", machine.Current.Name)
	msg := tea.KeyPressMsg{Text: "c", Mod: 0, Code: 99, ShiftedCode: 0, BaseCode: 0, IsRepeat: false}
	state, ok := machine.Next(msg)
	assert.True(t, ok)
	assert.Equal(t, "containers", state.Name)
}
