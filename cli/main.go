package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bcicen/grmon"
	ui "github.com/gizak/termui"
)

var (
	client      = &http.Client{Timeout: 10 * time.Second}
	grid        Grid
	lastRefresh time.Time
)

// parse command line arguments
var (
	helpFlag     = flag.Bool("h", false, "display this help dialog")
	hostFlag     = flag.String("host", "localhost:1234", "listening grmon host")
	endpointFlag = flag.String("endpoint", "/debug/grmon", "URL endpoint for grmon")
	intervalFlag = flag.Int("i", 0, "time in seconds between refresh")
)

type Grid []*widgets

func (g Grid) Len() int      { return len(g) }
func (g Grid) Swap(i, j int) { g[i], g[j] = g[j], g[i] }
func (g Grid) Less(i, j int) bool {
	return g[i].r.Num < g[j].r.Num
}

func (g Grid) Align() {
	sort.Sort(g)

	y := 2
	for _, w := range g {
		w.SetY(y)
		y += w.Height()
	}
}

func (g Grid) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, w := range g {
		buf.Merge(w.Buffer())
	}
	return buf
}

type widgets struct {
	r         grmon.Routine
	num       *ui.Par
	state     *ui.Par
	desc      *ui.Par
	trace     *ui.List
	showTrace bool
}

func (w *widgets) SetState(s string) {
	s = strings.Split(s, ",")[0]
	switch s {
	case "running":
		w.state.TextFgColor = ui.ColorGreen
	case "sleep":
		w.state.TextFgColor = ui.ColorYellow
	case "IO wait":
		w.state.TextFgColor = ui.ColorYellow
	default:
		w.state.TextFgColor = ui.ColorDefault
	}
	w.state.Text = s
}

func (w *widgets) SetTrace(a []string) {
	var lines []string
	for _, s := range a {
		lines = append(lines, strings.Replace(s, "\t", "  ", -1))
	}
	w.trace.Items = lines
	w.trace.Height = len(lines)
}

func (w *widgets) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(w.num.Buffer())
	buf.Merge(w.state.Buffer())
	buf.Merge(w.desc.Buffer())
	if w.showTrace {
		buf.Merge(w.trace.Buffer())
	}
	return buf
}

func (w *widgets) Height() int {
	if !w.showTrace {
		return 1
	}
	return len(w.trace.Items) + 1
}

func (w *widgets) SetY(y int) {
	w.num.Y = y
	w.state.Y = y
	w.desc.Y = y
	w.trace.Y = y + 1
}

func newWidgets() *widgets {
	p0 := ui.NewPar("")
	p0.X = 1
	p0.Height = 1
	p0.Width = 5
	p0.Border = false
	p0.TextFgColor = ui.ColorCyan

	p1 := ui.NewPar("")
	p1.X = 6
	p1.Height = 1
	p1.Width = 20
	p1.Border = false

	p2 := ui.NewPar("")
	p2.X = 26
	p2.Height = 1
	p2.Width = ui.TermWidth() - 26
	p2.Border = false

	ls := ui.NewList()
	ls.X = 2
	ls.Border = false
	ls.Width = ui.TermWidth()

	return &widgets{
		num:   p0,
		state: p1,
		desc:  p2,
		trace: ls,
	}
}

func Refresh() {

	bg := ui.NewPar("")
	bg.Width = ui.TermWidth()
	bg.Height = 1
	bg.Border = false
	bg.Y = 0
	bg.Bg = ui.ColorWhite

	header := newWidgets()
	header.SetY(0)
	header.num.Text = "#"
	header.num.Bg = ui.ColorWhite
	header.num.TextBgColor = ui.ColorWhite
	header.num.TextFgColor = ui.ColorBlack

	header.state.Text = "state"
	header.state.Bg = ui.ColorWhite
	header.state.TextBgColor = ui.ColorWhite
	header.state.TextFgColor = ui.ColorBlack

	header.desc.Text = "desc"
	header.desc.Bg = ui.ColorWhite
	header.desc.TextBgColor = ui.ColorWhite
	header.desc.TextFgColor = ui.ColorBlack

	grid = Grid{}

	for _, r := range poll() {
		w := newWidgets()
		w.r = r
		w.num.Text = fmt.Sprintf("%d", r.Num)
		w.SetState(r.State)
		w.desc.Text = r.Trace[0]

		r.Trace[0] = r.State
		w.SetTrace(r.Trace)
		grid = append(grid, w)
	}

	grid.Align()

	ui.Clear()
	ui.Render(bg)
	ui.Render(header)
	ui.Render(grid)
	lastRefresh = time.Now()
}

func main() {
	flag.Parse()

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	if *intervalFlag > 0 {
		ui.Handle("/timer/1s", func(ui.Event) {
			if int(time.Since(lastRefresh).Seconds()) > *intervalFlag {
				Refresh()
			}
		})
	}

	ui.Handle("/sys/kbd/r", func(ui.Event) {
		Refresh()
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	Refresh()
	ui.Loop()
}

var helpMsg = `grmon - goroutine monitor

usage: grmon [options]

options:
`

func poll() (routines []grmon.Routine) {
	url := fmt.Sprintf("http://%s%s", *hostFlag, *endpointFlag)
	r, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&routines)
	if err != nil {
		panic(err)
	}

	return
}

func printHelp() {
	fmt.Println(helpMsg)
	flag.PrintDefaults()
}
