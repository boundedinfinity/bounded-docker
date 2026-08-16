package widget

import (
	btable "charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/boundedinfinity/docker-tui/tui/menu"
	"github.com/boundedinfinity/docker-tui/tui/table"
)

func Widget[D any](title string, columns []btable.Column, data2RowFunc func(int, D) btable.Row, data []D) (tea.Model, tea.Model) {
	t := table.New(title, columns, data2RowFunc, data)
	m := menu.Item[D](title)

	return t, m
}
