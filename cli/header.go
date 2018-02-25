package main

import (
	"fmt"

	ui "github.com/gizak/termui"
)

var DefaultHeader = newHeader()

type header struct {
	*widgets
	bg *ui.Par
	ts *ui.Par
}

func (h *header) Align() {
	h.bg.Width = ui.TermWidth()
	h.ts.Width = ui.TermWidth()
	h.widgets.Align()
}

func (h *header) Buffer() ui.Buffer {
	t := "-"
	if !lastRefresh.IsZero() {
		t = lastRefresh.Format("15:04:05 MST")
	}
	h.ts.Text = fmt.Sprintf("last update: %s", t)

	buf := ui.NewBuffer()
	buf.Merge(h.bg.Buffer())
	buf.Merge(h.widgets.Buffer())
	buf.Merge(h.ts.Buffer())
	return buf
}

func newHeader() *header {
	bg := ui.NewPar("")
	bg.Height = 1
	bg.Border = false
	bg.Bg = ui.ColorWhite

	ts := ui.NewPar("-")
	ts.X = 2
	ts.Y = 1
	ts.Height = 1
	ts.Border = false
	ts.TextFgColor = ui.Attribute(248)

	h := newWidgets()
	h.num.Text = "#"
	h.num.Bg = ui.ColorWhite
	h.num.TextBgColor = ui.ColorWhite
	h.num.TextFgColor = ui.ColorBlack

	h.state.Text = "state"
	h.state.Bg = ui.ColorWhite
	h.state.TextBgColor = ui.ColorWhite
	h.state.TextFgColor = ui.ColorBlack

	h.desc.Text = "desc"
	h.desc.Bg = ui.ColorWhite
	h.desc.TextBgColor = ui.ColorWhite
	h.desc.TextFgColor = ui.ColorBlack

	return &header{h, bg, ts}
}
