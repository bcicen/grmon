package main

import (
	"fmt"

	ui "github.com/gizak/termui"
)

var DefaultHeader = newHeader()

type header struct {
	*widgets
	bg    *ui.Par
	ts    *ui.Par
	count *ui.Par
}

func (h *header) Align() {
	h.bg.Width = ui.TermWidth()

	h.ts.Y = ui.TermHeight() - 1
	h.ts.X = ui.TermWidth() - h.ts.Width

	h.count.Y = ui.TermHeight() - 1

	h.widgets.Align()
}

func (h *header) Update() {
	t := "-"
	if !lastRefresh.IsZero() {
		t = lastRefresh.Format("15:04:05 MST")
	}
	h.ts.Text = fmt.Sprintf("last update: %s", t)
	h.count.Text = fmt.Sprintf("total: %d", len(grid.rows))
}

func (h *header) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(h.bg.Buffer())
	buf.Merge(h.widgets.Buffer())
	buf.Merge(h.count.Buffer())
	buf.Merge(h.ts.Buffer())
	return buf
}

func newHeader() *header {
	bg := ui.NewPar("")
	bg.Height = 1
	bg.Border = false
	bg.Bg = ui.ColorWhite

	ts := ui.NewPar("-")
	ts.Height = 1
	ts.Width = 27
	ts.Border = false
	ts.TextFgColor = ui.Attribute(248)

	count := ui.NewPar("-")
	count.X = 1
	count.Height = 1
	count.Width = 40
	count.Border = false
	count.TextFgColor = ui.Attribute(248)

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

	return &header{h, bg, ts, count}
}
