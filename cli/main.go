package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	ui "github.com/gizak/termui"
	"github.com/nsf/termbox-go"
)

var (
	grid        = NewGrid()
	lastRefresh time.Time
)

// parse command line arguments
var (
	helpFlag     = flag.Bool("h", false, "display this help dialog")
	hostFlag     = flag.String("host", "localhost:1234", "listening grmon host")
	endpointFlag = flag.String("endpoint", "/debug/grmon", "URL endpoint for grmon")
	intervalFlag = flag.Int("i", 0, "time in seconds between refresh")
)

func Refresh() {
	routines, err := poll()
	if err == nil {
		grid.Clear()

		for _, r := range routines {
			w := newWidgets()
			w.num.Text = fmt.Sprintf("%d", r.Num)
			w.SetState(r.State)
			w.desc.Text = r.Trace[0]

			r.Trace[0] = r.State
			w.SetTrace(r.Trace)
			grid.AddRow(w)
		}
		lastRefresh = time.Now()
	}

	Render()
}

func Render() {
	grid.Align()
	ui.Clear()
	ui.Render(grid)
}

func HelpDialog() {
	p := ui.NewList()
	p.X = 1
	p.Height = 6
	p.Width = 45
	p.BorderLabel = "help"
	p.Items = []string{
		" r - refresh",
		" <up>,<down>,j,k - move cursor position",
		" <enter>,o - expand trace under cursor",
		" <esc>,q - exit grmon",
	}
	ui.Clear()
	ui.Render(p)
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func Display() bool {
	var menu func()

	if *intervalFlag > 0 {
		ui.Handle("/timer/1s", func(ui.Event) {
			if int(time.Since(lastRefresh).Seconds()) > *intervalFlag {
				Refresh()
			}
		})
	}

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		Render()
	})

	HandleKeys("up", func() {
		if grid.cursorPos > 0 {
			grid.cursorPos--
			Render()
		}
	})

	HandleKeys("down", func() {
		if grid.cursorPos < len(grid.rows)-1 {
			grid.cursorPos++
			Render()
		}
	})

	HandleKeys("enter", func() {
		grid.rows[grid.cursorPos].ToggleShowTrace()
		Render()
	})

	ui.Handle("/sys/kbd/r", func(ui.Event) {
		Refresh()
	})

	HandleKeys("exit", func() {
		ui.StopLoop()
	})

	HandleKeys("help", func() {
		menu = HelpDialog
		ui.StopLoop()
	})

	Render()
	ui.Loop()
	if menu != nil {
		menu()
		return false
	}
	return true
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
	termbox.SetOutputMode(termbox.Output256)

	Refresh()

	var quit bool
	for !quit {
		quit = Display()
	}
}

var helpMsg = `grmon - goroutine monitor

usage: grmon [options]

options:
`

func printHelp() {
	fmt.Println(helpMsg)
	flag.PrintDefaults()
}
