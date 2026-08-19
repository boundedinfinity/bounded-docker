package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/flosch/go-humanize"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
	"github.com/moby/moby/api/types/network"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////

var Utils = utils{}

type utils struct {
}

func (_ utils) unix2Time(unix int64) string {
	d := time.Unix(unix, 0)
	return d.Format(time.RFC3339)
}

func (_ utils) labels2Str(labels map[string]string) string {
	labelStrs := []string{}
	for k, v := range labels {
		labelStrs = append(labelStrs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(labelStrs, ",")
}

func (_ utils) ports2Str(ports []container.PortSummary) string {
	portStrs := []string{}
	for _, port := range ports {
		portStrs = append(portStrs, fmt.Sprintf("%s:%d->%d/%s", port.IP, port.PrivatePort, port.PublicPort, port.Type))
	}
	return strings.Join(portStrs, ",")
}

func (_ utils) index2Str(i int) string {
	return fmt.Sprintf("%d", i+1)
}

func (_ utils) size2Str(s int64) string {
	return humanize.Bytes(uint64(s))
}

func (_ utils) strNormal(text string, colWidth int) string {
	text = strings.TrimSpace(text)
	text = Utils.truncateString(text, "...", colWidth)
	return text
}

func (_ utils) truncateString(s, suffix string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}

	suffixWidth := len(suffix)

	if maxWidth <= suffixWidth {
		return s[:maxWidth]
	}

	return s[:maxWidth-suffixWidth] + suffix
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

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

// /////////////////////////////////////////////////////////////////////////////////////////////////

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

func imageTitles() []string {
	return []string{"#", "ID", "Image", "Size"}
}

func image2Row(i int, summary image.Summary) []string {
	return []string{
		Utils.index2Str(i),
		summary.ID,
		strings.Join(summary.RepoTags, "\n"),
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

// /////////////////////////////////////////////////////////////////////////////////////////////////

func createFakeErrors() []error {
	return []error{}
	// return []error{
	// 	errors.New("error 1"),
	// 	errors.New("error 2"),
	// 	errors.New("error 3"),
	// }
}

func errorTitles() []string {
	return []string{"#", "Error"}
}

func error2Row(i int, err error) []string {
	return []string{
		Utils.index2Str(i),
		err.Error(),
	}
}
