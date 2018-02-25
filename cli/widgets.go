package main

import (
	"strings"

	ui "github.com/gizak/termui"
)

type widgets struct {
	num       *ui.Par
	state     *ui.Par
	desc      *ui.Par
	trace     *ui.List
	y         int
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

func (w *widgets) ToggleShowTrace() { w.showTrace = w.showTrace != true }

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
	w.y = y
	w.num.Y = y
	w.state.Y = y
	w.desc.Y = y
	w.trace.Y = y + 1
}

func (w *widgets) Align() {
	w.desc.Width = ui.TermWidth() - 27
	w.trace.Width = ui.TermWidth()
}

func newWidgets() *widgets {
	p0 := ui.NewPar("")
	p0.X = 2
	p0.Height = 1
	p0.Width = 5
	p0.Border = false
	p0.TextBgColor = ui.ColorDefault

	p1 := ui.NewPar("")
	p1.X = 7
	p1.Height = 1
	p1.Width = 20
	p1.Border = false

	p2 := ui.NewPar("")
	p2.X = 27
	p2.Height = 1
	p2.Border = false

	ls := ui.NewList()
	ls.X = 7
	ls.Border = false

	return &widgets{
		num:   p0,
		state: p1,
		desc:  p2,
		trace: ls,
	}
}
