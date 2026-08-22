package tui

import (
	"github.com/moby/moby/api/types/image"
)

type imageUtils struct{}

func (_ imageUtils) titles() []string {
	return []string{"#", "ID", "Image", "Size"}
}

func (_ imageUtils) id(summary image.Summary) string {
	return summary.ID
}

func (_ imageUtils) summary2row(i int, summary image.Summary) []string {
	return []string{
		Utils.index2Str(i),
		summary.ID,
		Utils.Docker.repoTags2Str(summary.RepoTags),
		Utils.size2Str(summary.Size),
	}
}

func (_ imageUtils) fake() []image.Summary {
	return []image.Summary{}

	// return []image.Summary{
	// 	{
	// 		ID:          "id",
	// 		RepoTags:    []string{"tag1", "tag2"},
	// 		RepoDigests: []string{"digest1", "digest2"},
	// 		ParentID:    "parent_id",
	// 		Created:     0,
	// 		Size:        100,
	// 		SharedSize:  100,
	// 		Labels:      map[string]string{"label1": "value1", "label2": "value2"},
	// 		Containers:  10,
	// 	},
	// }
}
