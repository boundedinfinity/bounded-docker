package tui

import (
	"github.com/moby/moby/api/types/image"
)

func imageTitles() []string {
	return []string{"#", "ID", "Image", "Size"}
}

func image2Row(i int, summary image.Summary) []string {
	return []string{
		Utils.index2Str(i),
		summary.ID,
		Utils.repoTags2Str(summary.RepoTags),
		Utils.size2Str(summary.Size),
	}
}

func createFakeImages() []image.Summary {
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
