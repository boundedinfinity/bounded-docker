package main

import (
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/moby/moby/api/types/container"
)

func newContainersModel(summaries []container.Summary) tea.Model {
	data2row := func(_ int, summary container.Summary) table.Row {
		names := strings.Join(summary.Names, ", ")
		return table.Row{
			summary.ID,
			names,
			summary.Image,
			summary.Command,
			summary.Status,
		}
	}

	columns := []table.Column{
		{Title: "ID", Width: 30},
		{Title: "Name", Width: 30},
		{Title: "Image", Width: 30},
		{Title: "Command", Width: 30},
		{Title: "Status", Width: 30},
	}

	return newTableModel("Containers", columns, data2row, summaries)
}

func createFakeContainers() []container.Summary {
	return []container.Summary{}
}
