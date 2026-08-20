package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/flosch/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/moby/moby/api/types/container"
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

func (_ utils) repoTags2Str(tags []string) string {
	return strings.Join(tags, "\n")
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

func (_ utils) tcellEvent2Str(event *tcell.EventKey) string {
	key := string(event.Name())
	key = strings.Replace(key, "Rune[", "", 1)
	key = strings.Replace(key, "]", "", 1)
	key = strings.ToLower(key)

	return key
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
