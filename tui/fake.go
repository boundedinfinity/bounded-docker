package tui

import (
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
)

func item2Rows[T any](items []T, item2RowFunc func(int, T) []string) [][]string {
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = item2RowFunc(i, item)
	}
	return rows
}

func containerTitles() []string {
	return []string{"#", "ID", "Image", "Command", "Created", "Status", "Names", "Ports", "Labels"}
}

func container2Row(i int, summary container.Summary) []string {
	labels := []string{}
	for k, v := range summary.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", k, v))
	}

	return []string{
		fmt.Sprintf("%3d", i+1),
		summary.ID,
		summary.Image,
		summary.Command,
		fmt.Sprintf("%d", summary.Created),
		summary.Status,
		strings.Join(summary.Names, "\n"),
		fmt.Sprintf("%v", summary.Ports),
		strings.Join(labels, "\n"),
	}
}

func createFakeContainers() []container.Summary {
	return []container.Summary{
		{
			ID:      "id",
			Image:   "image",
			Command: "command",
			Created: 0,
			State:   "state",
			Status:  "status",
			Ports: []container.PortSummary{
				{IP: netip.Addr{}, PrivatePort: 80, PublicPort: 80, Type: "type"},
			},
			Labels:     map[string]string{"label1": "value1", "label2": "value2"},
			SizeRw:     100,
			SizeRootFs: 100,
			Mounts: []container.MountPoint{
				{
					Type:        "type",
					Name:        "name",
					Source:      "source",
					Destination: "destination",
					Driver:      "driver",
					Mode:        "mode",
					RW:          true,
					Propagation: "propagation",
				},
			},
			Names: []string{"name1", "name2"},
		},
	}
}

func imageTitles() []string {
	return []string{"#", "ID", "Image", "Size"}
}

func image2Row(i int, summary image.Summary) []string {
	return []string{
		fmt.Sprintf("%3d", i+1),
		summary.ID,
		strings.Join(summary.RepoTags, "\n"),
		fmt.Sprintf("%d", summary.Size),
	}
}

func createFakeImages() []image.Summary {
	return []image.Summary{
		{
			ID:          "id",
			RepoTags:    []string{"tag1", "tag2"},
			RepoDigests: []string{"digest1", "digest2"},
			ParentID:    "parent_id",
			Created:     0,
			Size:        100,
			SharedSize:  100,
			Labels:      map[string]string{"label1": "value1", "label2": "value2"},
			Containers:  10,
		},
	}
}

func createFakeErrors() []error {
	return []error{
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("error 3"),
	}
}

func errorTitles() []string {
	return []string{"#", "Error"}
}

func error2Row(i int, err error) []string {
	return []string{
		fmt.Sprintf("%3d", i+1),
		err.Error(),
	}
}
