package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	ui "github.com/gizak/termui"
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
	grid.Clear()

	routines := poll()
	for _, r := range routines {
		w := newWidgets()
		w.num.Text = fmt.Sprintf("%d", r.Num)
		w.SetState(r.State)
		w.desc.Text = r.Trace[0]

		r.Trace[0] = r.State
		w.SetTrace(r.Trace)
		grid.AddRow(w)
	}

	Render()
	lastRefresh = time.Now()
}

func Render() {
	grid.Align()
	ui.Clear()
	ui.Render(grid)
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

func printHelp() {
	fmt.Println(helpMsg)
	flag.PrintDefaults()
}
