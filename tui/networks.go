package tui

import (
	"github.com/moby/moby/api/types/network"
)

func networkTitles() []string {
	return []string{"#", "ID", "Name", "Driver", "Scope", "Labels"}
}

func network2Row(i int, summary network.Summary) []string {
	return []string{
		Utils.index2Str(i),
		summary.ID,
		summary.Name,
		summary.Driver,
		summary.Scope,
		Utils.labels2Str(summary.Labels),
	}
}
