package tui

import (
	"fmt"
	"math"
	"net/netip"
	"sort"
	"strings"
	"time"

	"github.com/boundedinfinity/docker-tui/state"
	"github.com/flosch/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
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
	Inspect   inspectUtils
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type inspectUtils struct{}

func (_ inspectUtils) env2Str(env []string) string {
	return strings.Join(env, "\n")
}

// flatten turns decoded JSON into sorted key/value pairs using dotted paths.
func (this inspectUtils) flatten(prefix string, value any, rows *[][2]string) {
	switch v := value.(type) {
	case map[string]any:
		if len(v) == 0 {
			*rows = append(*rows, [2]string{prefix, "{}"})
			return
		}

		keys := make([]string, 0, len(v))

		for k := range v {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			this.flatten(this.join(prefix, k), v[k], rows)
		}
	case []any:
		if len(v) == 0 {
			*rows = append(*rows, [2]string{prefix, "[]"})
			return
		}

		for i, item := range v {
			this.flatten(fmt.Sprintf("%s[%d]", prefix, i), item, rows)
		}
	case nil:
		*rows = append(*rows, [2]string{prefix, ""})
	default:
		*rows = append(*rows, [2]string{prefix, fmt.Sprint(v)})
	}
}

func (_ inspectUtils) join(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return prefix + "." + key
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

func (_ dockerUtils) map2Lines(values map[string]string) []string {
	lines := make([]string, 0, len(values))

	for k, v := range values {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(lines)
	return lines
}

func (_ dockerUtils) volumes2Lines(volumes map[string]struct{}) []string {
	lines := make([]string, 0, len(volumes))

	for k := range volumes {
		lines = append(lines, k)
	}

	sort.Strings(lines)
	return lines
}

func (_ dockerUtils) addrs2Lines(addrs []netip.Addr) []string {
	lines := make([]string, 0, len(addrs))

	for _, addr := range addrs {
		lines = append(lines, addr.String())
	}

	return lines
}

func (_ dockerUtils) exposedPorts2Lines(ports network.PortSet) []string {
	lines := make([]string, 0, len(ports))

	for port := range ports {
		lines = append(lines, port.String())
	}

	sort.Strings(lines)
	return lines
}

func (_ dockerUtils) portMap2Lines(ports network.PortMap) []string {
	lines := make([]string, 0, len(ports))

	for port, bindings := range ports {
		if len(bindings) == 0 {
			lines = append(lines, port.String())
			continue
		}

		for _, binding := range bindings {
			lines = append(lines, fmt.Sprintf("%s:%s->%s", binding.HostIP, binding.HostPort, port))
		}
	}

	sort.Strings(lines)
	return lines
}

// /////////////////////////////////////////////////////////////////////////////////////////////////

type stringUtils struct{}

func (this stringUtils) fixed(text string, width int, suffix string) string {
	if len(text) > width {
		text = this.truncateString(text, width, suffix)
	}

	if len(text) < width {
		text = this.pad(text, width)
	}

	return text
}

func (this stringUtils) fixedFn(width int, suffix string) func(string) string {
	return func(text string) string { return this.fixed(text, width, suffix) }
}

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
	text = Utils.String.truncateString(text, colWidth, "...")
	return text
}

func (this stringUtils) clamp(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}

	return this.truncateString(text, maxWidth, "...")
}

func (_ stringUtils) truncateString(text string, maxWidth int, suffix string) string {
	if len(text) <= maxWidth {
		return text
	}

	suffixWidth := len(suffix)

	if maxWidth <= suffixWidth {
		return text[:maxWidth]
	}

	return text[:maxWidth-suffixWidth] + suffix
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

func (_ tviewUtils) multiRowi(rows [][]string, index int, title string, lines ...string) ([][]string, int) {
	next := func() string {
		s := Utils.index2Str(index)
		index++
		return s
	}

	if len(lines) == 0 {
		return append(rows, []string{next(), title, ""}), index
	}

	rows = append(rows, []string{next(), title, lines[0]})
	for _, line := range lines[1:] {
		rows = append(rows, []string{next(), "", line})
	}

	return rows, index
}

func (_ tviewUtils) multiRow(rows [][]string, title string, lines ...string) [][]string {
	if len(lines) == 0 {
		return append(rows, []string{title, ""})
	}

	rows = append(rows, []string{title, lines[0]})
	for _, line := range lines[1:] {
		rows = append(rows, []string{"", line})
	}

	return rows
}

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

func (_ errorUtils) error2rows(errs []error) [][]string {
	rows := [][]string{}
	rows = append(rows, []string{"#", "Error"})

	for i, err := range errs {
		rows = append(rows, []string{
			Utils.index2Str(i),
			err.Error(),
		})
	}

	return rows
}
