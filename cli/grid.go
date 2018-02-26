package main

import (
	ui "github.com/gizak/termui"
)

type Grid struct {
	header    *header
	rows      []*widgets
	cursor    *ui.Par
	cursorPos int
}

func NewGrid() *Grid {
	c := ui.NewPar(">")
	c.Width = 1
	c.Border = false
	c.TextFgColor = ui.ColorCyan
	return &Grid{
		header: newHeader(),
		cursor: c,
	}
}

func (g *Grid) AddRow(w *widgets) { g.rows = append(g.rows, w) }
func (g *Grid) Clear() {
	g.cursorPos = 0
	g.rows = []*widgets{}
}

func (g *Grid) Align() {
	g.header.Align()

	y := 2
	for n, w := range g.rows {
		w.SetY(y)
		w.Align()
		y += w.Height()
		if n == g.cursorPos {
			g.cursor.Y = w.y
		}
	}
}

func (g *Grid) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(g.header.Buffer())
	for _, w := range g.rows {
		buf.Merge(w.Buffer())
	}
	if len(g.rows) > 0 {
		buf.Merge(g.cursor.Buffer())
	}
	return buf
}
