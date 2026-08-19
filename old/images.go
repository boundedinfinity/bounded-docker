package main

// func newImageModel(summaries []image.Summary) (tea.Model, tea.Model) {
// 	data2row := func(i int, summary image.Summary) table.Row {
// 		var image string

// 		if len(summary.RepoTags) > 0 {
// 			image = summary.RepoTags[0]
// 		} else {
// 			image = "<none>"
// 		}

// 		size := humanize.Bytes(uint64(summary.Size))

// 		return table.Row{
// 			fmt.Sprintf("%3d", i+1),
// 			summary.ID,
// 			image,
// 			size,
// 		}
// 	}

// 	columns := []table.Column{
// 		{Title: "#", Width: 3},
// 		{Title: "ID", Width: 30},
// 		{Title: "Image", Width: 30},
// 		{Title: "Size", Width: 30},
// 	}

// 	return widget.Widget("images", columns, data2row, summaries)
// }

// func createFakeImages() []image.Summary {
// 	return []image.Summary{}
// }
