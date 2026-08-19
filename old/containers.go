package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/moby/moby/api/types/container"
)

func newContainersModel(summaries []container.Summary) (tea.Model, tea.Model) {
	// data2row := func(i int, summary container.Summary) table.Row {
	// 	names := strings.Join(summary.Names, ", ")

	// 	return table.Row{
	// 		fmt.Sprintf("%3d", i+1),
	// 		summary.ID,
	// 		names,
	// 		summary.Image,
	// 		summary.Command,
	// 		summary.Status,
	// 	}
	// }

	// // data2details := func(i int, summary container.Summary) []table.Row {
	// // 	return []table.Row{
	// // 		{"#", fmt.Sprintf("%3d", i+1)},
	// // 		{"ID", summary.ID},
	// // 		{"Image", summary.Image},
	// // 		{"Command", summary.Command},
	// // 		{"Created", fmt.Sprintf("%d", summary.Created)},
	// // 		{"Status", summary.Status},
	// // 		{"Names", strings.Join(summary.Names, "\n")},
	// // 		{"Ports", fmt.Sprintf("%v", summary.Ports)},
	// // 		{"Labels", fmt.Sprintf("%v", summary.Labels)},
	// // 	}
	// // }

	// columns := []table.Column{
	// 	{Title: "#", Width: 3},
	// 	{Title: "ID", Width: 30},
	// 	{Title: "Name", Width: 30},
	// 	{Title: "Image", Width: 30},
	// 	{Title: "Command", Width: 30},
	// 	{Title: "Status", Width: 30},
	// }

	// return widget.Widget("containers", columns, data2row, summaries)
	return nil, nil
}
