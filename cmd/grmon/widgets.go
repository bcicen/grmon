package main

import (
	"strconv"
	"strings"
	"sync"

	ui "github.com/bcicen/termui"
)

type WidgetMap struct {
	m    map[int]*widgets
	lock sync.RWMutex
}

func NewWidgetMap() *WidgetMap {
	return &WidgetMap{m: make(map[int]*widgets)}
}

func (wm *WidgetMap) MustGet(id int) *widgets {
	if w, ok := wm.m[id]; ok {
		return w
	}
	return wm.add(id)
}

func (wm *WidgetMap) Del(id int) {
	wm.lock.Lock()
	defer wm.lock.Unlock()
	if _, ok := wm.m[id]; ok {
		delete(wm.m, id)
	}
}

func (wm *WidgetMap) add(id int) *widgets {
	wm.lock.Lock()
	defer wm.lock.Unlock()
	w := newWidgets()
	w.num.Text = strconv.Itoa(id)
	w.Align()
	wm.m[id] = w
	return w
}

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

type Column interface {
	ui.Bufferer
	SetY(int)
}

func (w *widgets) cols() []Column {
	a := []Column{
		w.num,
		w.state,
		w.desc,
	}
	return a
}

func (w *widgets) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, b := range w.cols() {
		buf.Merge(b.Buffer())
	}
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
	if y != w.y {
		w.y = y
		for _, b := range w.cols() {
			b.SetY(y)
		}
	}
	w.trace.Y = y + 1
}

func (w *widgets) Align() {
	w.desc.Width = ui.TermWidth() - 27
	w.trace.Width = ui.TermWidth()
}

func newWidgets() *widgets {
	num := newCol(2, 5)
	state := newCol(7, 20)
	desc := newCol(27, 20)

	trace := ui.NewList()
	trace.X = 7
	trace.Border = false

	return &widgets{
		num:   num,
		state: state,
		desc:  desc,
		trace: trace,
	}
}

func newCol(x, w int) *ui.Par {
	p := ui.NewPar("")
	p.X = x
	p.Height = 1
	p.Width = w
	p.Border = false
	p.TextBgColor = ui.ColorDefault
	return p
}
