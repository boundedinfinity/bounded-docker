package tui

import (
	"strings"

	"github.com/moby/moby/api/types/container"
)

type containerUtils struct{}

func (_ containerUtils) id(summary container.Summary) string {
	return summary.ID
}

func (_ containerUtils) summary2rows(summaries []container.Summary) [][]string {
	rows := [][]string{}
	rows = append(rows, []string{"#", "ID", "Image", "Command", "Created", "Status", "Names", "Ports", "Labels"})

	for i, summary := range summaries {
		rows = append(rows, []string{
			Utils.index2Str(i),
			summary.ID,
			summary.Image,
			summary.Command,
			Utils.Time.unix2Time(summary.Created),
			summary.Status,
			strings.Join(summary.Names, ","),
			Utils.Docker.ports2Str(summary.Ports),
			Utils.Docker.labels2Str(summary.Labels),
		})
	}

	return rows
}
