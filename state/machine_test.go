package state_test

import (
	"testing"

	"github.com/boundedinfinity/docker-tui/state"
	"github.com/stretchr/testify/assert"
)

// v2.KeyPressMsg {Text: "c", Mod: 0, Code: 99, ShiftedCode: 0, BaseCode: 0, IsRepeat: false}
// v2.KeyPressMsg {Text: "e", Mod: 0, Code: 101, ShiftedCode: 0, BaseCode: 0, IsRepeat: false}
// v2.KeyPressMsg {Text: "i", Mod: 0, Code: 105, ShiftedCode: 0, BaseCode: 0, IsRepeat: false}

func Test_StateMachine_Creation_Default(t *testing.T) {
	machine, err := state.New(state.DefaultConfig())

	assert.NoError(t, err)
	assert.Equal(t, "root", machine.Current.Id)

	// c := tea.KeyPressMsg{Text: "c", Mod: 0, Code: 99, ShiftedCode: 0, BaseCode: 0, IsRepeat: false}
	next, ok := machine.Next("c")
	assert.True(t, ok)
	assert.Equal(t, "containers", next.Id)

	// up := tea.KeyPressMsg{Text: "j", Mod: 0, Code: 106, ShiftedCode: 0, BaseCode: 0, IsRepeat: false}
	next, ok = machine.Next("up")
	assert.True(t, ok)
	assert.Equal(t, "containers", next.Id)
}
