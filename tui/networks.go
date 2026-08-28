package tui

import (
	"github.com/moby/moby/api/types/network"
)

type networkUtils struct{}

func (_ networkUtils) id(summary network.Summary) string {
	return summary.ID
}

func (_ networkUtils) summary2rows(summaries []network.Summary) [][]string {
	rows := [][]string{}
	rows = append(rows, []string{"#", "ID", "Name", "Driver", "Scope", "Labels"})

	for i, summary := range summaries {
		rows = append(rows, []string{
			Utils.index2Str(i),
			summary.ID,
			summary.Name,
			summary.Driver,
			summary.Scope,
			Utils.Docker.labels2Str(summary.Labels),
		})
	}

	return rows
}
