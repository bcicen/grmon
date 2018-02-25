package main

import (
	ui "github.com/gizak/termui"
)

var DefaultHeader = newHeader()

type header struct {
	*widgets
	bg *ui.Par
}

func (h *header) Align() {
	h.bg.Width = ui.TermWidth()
	h.widgets.Align()
}

func (h *header) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(h.bg.Buffer())
	buf.Merge(h.widgets.Buffer())
	return buf
}

func newHeader() *header {
	bg := ui.NewPar("")
	bg.Height = 1
	bg.Border = false
	bg.Bg = ui.ColorWhite

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

	return &header{h, bg}
}
