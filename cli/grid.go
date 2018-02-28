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
	maxRows   int
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

func (g *Grid) CursorUp() (ok bool) {
	if g.cursorPos > 0 {
		g.cursorPos--
		nextRowHeight := g.rows[g.cursorPos].Height()
		if g.cursor.Y-nextRowHeight < 2 {
			g.y += nextRowHeight
		}
		ok = true
	}

	return
}

func (g *Grid) CursorDown() (ok bool) {
	if g.cursorPos < len(g.rows)-1 {
		// if currently select row is beyond lower boundary
		// shift page up
		curRowHeight := g.rows[g.cursorPos].Height()
		if curRowHeight+g.cursor.Y-2 >= g.maxRows {
			g.y -= curRowHeight
		}
		g.cursorPos++
		ok = true
	}
	return
}

func (g *Grid) Clear() {
	g.cursorPos = 0
	g.rows = []*widgets{}
}

func (g *Grid) Align() {
	g.header.Align()
	for _, w := range g.rows {
		w.Align()
	}
	g.maxRows = ui.TermHeight() - 4
}

func (g *Grid) Buffer() ui.Buffer {
	buf := ui.NewBuffer()

	y := g.y

	for n, w := range g.rows {
		w.SetY(y)
		buf.Merge(w.Buffer())
		if paused && n == g.cursorPos {
			g.cursor.Y = y
			buf.Merge(g.cursor.Buffer())
		}
		y += w.Height()
	}

	buf.Merge(g.header.Buffer())
	return buf
}
