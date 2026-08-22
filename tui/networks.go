package tui

import (
	"github.com/moby/moby/api/types/network"
)

type networkUtils struct{}

func (_ networkUtils) titles() []string {
	return []string{"#", "ID", "Name", "Driver", "Scope", "Labels"}
}

func (_ networkUtils) id(summary network.Summary) string {
	return summary.ID
}

func (_ networkUtils) summary2rows(i int, summary network.Summary) []string {
	return []string{
		Utils.index2Str(i),
		summary.ID,
		summary.Name,
		summary.Driver,
		summary.Scope,
		Utils.Docker.labels2Str(summary.Labels),
	}
}
