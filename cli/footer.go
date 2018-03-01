package main

import (
	"fmt"

	ui "github.com/gizak/termui"
)

const footerHeight = 2

type footer struct {
	bg    *ui.Par
	ts    *ui.Par
	count *ui.Par
}

func (f *footer) Align() {
	f.SetY(ui.TermHeight() - footerHeight)
	f.bg.Width = ui.TermWidth()
	f.ts.X = ui.TermWidth() - f.ts.Width
}

func (f *footer) SetY(y int) {
	f.bg.Y = y
	y += (footerHeight - 1)
	f.ts.Y = y
	f.count.Y = y
}

func (f *footer) Update() {
	t := "-"
	if !lastRefresh.IsZero() {
		t = lastRefresh.Format("15:04:05 MST")
	}
	f.ts.Text = fmt.Sprintf("last update: %s", t)
	f.count.Text = fmt.Sprintf("total: %d", len(grid.rows))
}

func (f *footer) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(f.bg.Buffer())
	buf.Merge(f.count.Buffer())
	buf.Merge(f.ts.Buffer())
	return buf
}

func newFooter() *footer {
	bg := ui.NewPar("")
	bg.Height = 2
	bg.Border = false
	//bg.Bg = ui.ColorDefault

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

	return &footer{bg, ts, count}
}
