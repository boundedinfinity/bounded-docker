package main

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/dustin/go-humanize"
	"github.com/moby/moby/api/types/image"
)

func newImageModel(summaries []image.Summary) tea.Model {
	data2row := func(_ int, summary image.Summary) table.Row {
		// style := lipgloss.NewStyle().Padding(0, 10).Align(lipgloss.Right)
		var image string

		if len(summary.RepoTags) > 0 {
			image = summary.RepoTags[0]
		} else {
			image = "<none>"
		}

		size := humanize.Bytes(uint64(summary.Size))
		// size = style.Render(size)

		return table.Row{
			summary.ID,
			image,
			size,
		}
	}

	columns := []table.Column{
		{Title: "ID", Width: 30},
		{Title: "Image", Width: 30},
		{Title: "Size", Width: 30},
	}

	return newTableModel("Images", columns, data2row, summaries)
}

func createFakeImages() []image.Summary {
	return []image.Summary{}
}
