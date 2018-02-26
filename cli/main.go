package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	ui "github.com/gizak/termui"
	"github.com/nsf/termbox-go"
)

var (
	wmap        = NewWidgetMap()
	grid        = NewGrid()
	ctx         = context.Background()
	paused      bool
	lastRefresh time.Time
)

// parse command line arguments
var (
	helpFlag     = flag.Bool("h", false, "display this help dialog")
	hostFlag     = flag.String("host", "localhost:1234", "listening grmon host")
	selfFlag     = flag.Bool("self", false, "monitor grmon itself")
	endpointFlag = flag.String("endpoint", "/debug/grmon", "URL endpoint for grmon")
	intervalFlag = flag.Int("i", 5, "time in seconds between refresh")
)

func Refresh() {
	routines, err := poll()
	if err == nil {
		grid.Clear()

		for _, r := range routines {
			w := wmap.MustGet(r.Num)
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
	grid.header.Update()
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
		" r - manual refresh",
		" p - pause/unpause automatic updates",
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
	var next func()

	rctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if *intervalFlag > 0 && !paused {
		go func() {
			for {
				select {
				case <-rctx.Done():
					return
				case <-time.After(500 * time.Millisecond):
					if int(time.Since(lastRefresh).Seconds()) >= *intervalFlag {
						Refresh()
					}
				}
			}
		}()
	}

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		grid.Align()
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

	ui.Handle("/sys/kbd/p", func(ui.Event) {
		next = func() { paused = paused != true }
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/r", func(ui.Event) {
		Refresh()
	})

	HandleKeys("exit", func() {
		ui.StopLoop()
	})

	HandleKeys("help", func() {
		next = HelpDialog
		ui.StopLoop()
	})

	grid.Align()
	Render()
	ui.Loop()
	if next != nil {
		next()
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

	if *selfFlag {
		go http.ListenAndServe(":1234", nil)
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
