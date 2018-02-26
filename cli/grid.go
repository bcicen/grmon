package main

import (
	ui "github.com/gizak/termui"
)

type Grid struct {
	header    *header
	rows      []*widgets
	cursor    *ui.Par
	cursorPos int
	y         int
	maxY      int
}

func NewGrid() *Grid {
	c := ui.NewPar(">")
	c.Width = 1
	c.Border = false
	c.TextFgColor = ui.ColorCyan
	return &Grid{
		header: newHeader(),
		cursor: c,
		y:      2,
	}
}

func (g *Grid) AddRow(w *widgets) { g.rows = append(g.rows, w) }

func (g *Grid) Clear() {
	g.cursorPos = 0
	g.rows = []*widgets{}
}

func (g *Grid) Align() {
	g.header.Align()
	for _, w := range g.rows {
		w.Align()
	}
	g.maxY = ui.TermHeight() - g.y
}

func (g *Grid) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(g.header.Buffer())

	y := g.y
	for n, w := range g.rows {
		w.SetY(y)
		buf.Merge(w.Buffer())
		if paused && n == g.cursorPos {
			g.cursor.Y = y
			buf.Merge(g.cursor.Buffer())
		}
		y += w.Height()
		if y >= g.maxY {
			break
		}
	}

	return buf
}
