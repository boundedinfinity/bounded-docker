package tui

import (
	"fmt"
	"strings"

	"github.com/moby/moby/api/types/container"
)

type containerUtils struct{}

func (_ containerUtils) titles() []string {
	return []string{"#", "ID", "Image", "Command", "Created", "Status", "Names", "Ports", "Labels"}
}

func (_ containerUtils) id(summary container.Summary) string {
	return summary.ID
}

func (_ containerUtils) summary2rows(i int, summary container.Summary) []string {
	return []string{
		Utils.index2Str(i),
		summary.ID,
		summary.Image,
		summary.Command,
		Utils.Time.unix2Time(summary.Created),
		summary.Status,
		strings.Join(summary.Names, ","),
		Utils.Docker.ports2Str(summary.Ports),
		Utils.Docker.labels2Str(summary.Labels),
	}
}

func (_ containerUtils) inspect2rows(i int, result container.InspectResponse) [][]string {
	rows := [][]string{
		{"ID", result.ID},
		{"Path", result.Path},
		{"Created", Utils.Time.strTime(result.Created)},
		{"Args", strings.Join(result.Args, ",")},
		{"State", string(result.State.Status)},
		{"", fmt.Sprintf("%d", result.State.Pid)},
		{"Image", result.Image},
		{"Name", result.Name},
	}

	return rows
}

func (_ containerUtils) fake() []container.Summary {
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
