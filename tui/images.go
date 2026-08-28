package tui

import (
	"github.com/moby/moby/api/types/image"
)

type imageUtils struct{}

func (_ imageUtils) id(summary image.Summary) string {
	return summary.ID
}

func (_ imageUtils) summary2rows(summaries []image.Summary) [][]string {
	rows := [][]string{}
	rows = append(rows, []string{"#", "ID", "Image", "Size"})

	for i, summary := range summaries {
		rows = append(rows, []string{
			Utils.index2Str(i),
			summary.ID,
			Utils.Docker.repoTags2Str(summary.RepoTags),
			Utils.size2Str(summary.Size),
		})
	}

	return rows
}
