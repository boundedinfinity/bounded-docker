package tui

import (
	"strings"

	"github.com/moby/moby/api/types/container"
)

func containerTitles() []string {
	return []string{"#", "ID", "Image", "Command", "Created", "Status", "Names", "Ports", "Labels"}
}

func container2Row(i int, summary container.Summary) []string {
	return []string{
		Utils.index2Str(i),
		summary.ID,
		summary.Image,
		summary.Command,
		Utils.unix2Time(summary.Created),
		summary.Status,
		strings.Join(summary.Names, ","),
		Utils.ports2Str(summary.Ports),
		Utils.labels2Str(summary.Labels),
	}
}

func createFakeContainers() []container.Summary {
	return []container.Summary{}

	// return []container.Summary{
	// 	{
	// 		ID:      "id",
	// 		Image:   "image",
	// 		Command: "command",
	// 		Created: 0,
	// 		State:   "state",
	// 		Status:  "status",
	// 		Ports: []container.PortSummary{
	// 			{IP: netip.Addr{}, PrivatePort: 80, PublicPort: 80, Type: "type"},
	// 		},
	// 		Labels:     map[string]string{"label1": "value1", "label2": "value2"},
	// 		SizeRw:     100,
	// 		SizeRootFs: 100,
	// 		Mounts: []container.MountPoint{
	// 			{
	// 				Type:        "type",
	// 				Name:        "name",
	// 				Source:      "source",
	// 				Destination: "destination",
	// 				Driver:      "driver",
	// 				Mode:        "mode",
	// 				RW:          true,
	// 				Propagation: "propagation",
	// 			},
	// 		},
	// 		Names: []string{"name1", "name2"},
	// 	},
	// }
}
