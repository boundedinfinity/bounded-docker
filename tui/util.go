package tui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/boundedinfinity/docker-tui/state"
	"github.com/flosch/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/moby/moby/api/types/container"
	"github.com/rivo/tview"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////

var Utils = utils{}

type utils struct {
	Time      timeUtils
	Docker    dockerUtils
	String    stringUtils
	Tview     tviewUtils
	Error     errorUtils
	Fake      fakeUtils
	Container containerUtils
	Network   networkUtils
	Image     imageUtils
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type timeUtils struct{}

func (_ timeUtils) unix2Time(unix int64) string {
	timestamp := time.Unix(unix, 0)
	return timestamp.Format(time.RFC3339)
}

func (_ timeUtils) strTime(unix string) string {
	timestamp, err := time.Parse(time.RFC3339, unix)

	if err != nil {
		panic(err)
	}

	return timestamp.Format(time.RFC3339)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type dockerUtils struct{}

func (_ dockerUtils) labels2Str(labels map[string]string) string {
	labelStrs := []string{}
	for k, v := range labels {
		labelStrs = append(labelStrs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(labelStrs, ",")
}

func (_ dockerUtils) ports2Str(ports []container.PortSummary) string {
	portStrs := []string{}
	for _, port := range ports {
		portStrs = append(portStrs, fmt.Sprintf("%s:%d->%d/%s", port.IP, port.PrivatePort, port.PublicPort, port.Type))
	}
	return strings.Join(portStrs, ",")
}

func (_ dockerUtils) repoTags2Str(tags []string) string {
	return strings.Join(tags, "\n")
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type stringUtils struct{}

func (_ stringUtils) padRight(text string, width int) string {
	width = max(width, len(text))
	size := width - len(text)
	padding := strings.Repeat(" ", size)
	return text + padding
}

func (_ stringUtils) padLeft(text string, width int) string {
	width = max(width, len(text))
	size := width - len(text)
	padding := strings.Repeat(" ", size)
	return padding + text
}

func (_ stringUtils) pad(text string, width int) string {
	width = max(width, len(text))
	size := width - len(text)
	size /= 2
	padding := strings.Repeat(" ", size)
	return padding + text + padding
}

func (_ stringUtils) strNormal(text string, colWidth int) string {
	text = strings.TrimSpace(text)
	text = Utils.String.truncateString(text, "...", colWidth)
	return text
}

func (_ stringUtils) truncateString(s, suffix string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}

	suffixWidth := len(suffix)

	if maxWidth <= suffixWidth {
		return s[:maxWidth]
	}

	return s[:maxWidth-suffixWidth] + suffix
}

func (_ stringUtils) calcMax(rows [][]string) int {
	width := 0

	for _, row := range rows {
		for _, text := range row {
			width = max(width, len(text))
		}
	}

	return width
}

func (_ stringUtils) calcMin(rows [][]string) int {
	width := math.MaxInt

	for _, row := range rows {
		for _, text := range row {
			width = min(width, len(text))
		}
	}

	width += 10
	return width
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

func (_ utils) size2Str(s int64) string {
	return humanize.Bytes(uint64(s))
}

func (_ utils) index2Str(i int) string {
	return fmt.Sprintf("%d", i+1)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type tviewUtils struct{}

func (_ tviewUtils) tcellEvent2Str(event *tcell.EventKey) string {
	key := string(event.Name())
	key = strings.Replace(key, "Rune[", "", 1)
	key = strings.Replace(key, "]", "", 1)
	key = strings.ToLower(key)

	return key
}

func (_ tviewUtils) state2Rows(state *state.State) [][]string {
	rows := [][]string{}
	row := []string{"Transistions:"}

	for _, transistion := range state.Transitions {
		var keys []string

		for _, key := range transistion.Keys {
			keys = append(keys, key.Name)
		}

		col1 := fmt.Sprintf("%s", strings.Join(keys, "/"))
		col2 := fmt.Sprintf("› %s", transistion.State.Name)
		row = append(row, col1, col2)
	}

	rows = append(rows, row)
	row = []string{"    Commands:"}

	for _, command := range state.Commands {
		var keys []string

		for _, key := range command.Keys() {
			keys = append(keys, key.Name)
		}

		col1 := fmt.Sprintf("%s", strings.Join(keys, "/"))
		col2 := fmt.Sprintf("› %s", command.Command.Name)
		row = append(row, col1, col2)
	}

	rows = append(rows, row)
	return rows
}

func (_ tviewUtils) state2Table(state *state.State) *tview.Table {
	rows := Utils.Tview.state2Rows(state)
	width := Utils.String.calcMax(rows) / 2

	table := tview.NewTable()
	table.SetBorderPadding(0, 0, 1, 1)

	for r, row := range rows {
		for c, text := range row {
			if c%2 == 0 {
				text = Utils.String.padRight(text, width)
			} else {
				text = Utils.String.padLeft(text, width)
			}

			cell := tview.NewTableCell(text).
				SetTextColor(tcell.ColorWhite).
				SetSelectable(false)

			if c%2 == 0 {
				cell.SetAlign(tview.AlignLeft)
			} else {
				cell.SetTextColor(tcell.ColorYellow)
				cell.SetAlign(tview.AlignRight)
			}

			table.SetCell(r, c, cell)
		}
	}

	return table
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type fakeUtils struct{}

func (_ fakeUtils) createFakeErrors() []error {
	return []error{}
	// return []error{
	// 	errors.New("error 1"),
	// 	errors.New("error 2"),
	// 	errors.New("error 3"),
	// }
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type errorUtils struct{}

func (_ errorUtils) errorTitles() []string {
	return []string{"#", "Error"}
}

func (_ errorUtils) id(err error) string {
	return ""
}

func (_ errorUtils) error2Row(i int, err error) []string {
	return []string{
		Utils.index2Str(i),
		err.Error(),
	}
}
