package main

import (
	ui "github.com/gizak/termui"
)

func TraceDialog() {
	p := ui.NewList()
	p.X = 1
	p.Height = ui.TermHeight()
	p.Width = ui.TermWidth()
	p.Border = false
	p.Items = grid.rows[grid.cursorPos].trace.Items
	ui.Clear()
	ui.Render(p)
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func HelpDialog() {
	p := ui.NewList()
	p.X = 1
	p.Height = 6
	p.Width = 45
	p.BorderLabel = "help"
	p.Items = []string{
		" r - manual refresh",
		" s - toggle sort column and refresh",
		" p - pause/unpause automatic updates",
		" <up>,<down>,j,k - move cursor position",
		" <enter>,o - expand trace under cursor",
		" t - open trace in full screen",
		" <esc>,q - exit grmon",
	}
	ui.Clear()
	ui.Render(p)
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func FilterDialog() {
	ui.ResetHandlers()
	defer ui.ResetHandlers()

	i := NewInput()
	i.BorderLabel = "Filter"
	i.SetY(ui.TermHeight() - i.Height)
	i.Data = filter
	ui.Render(i)

	// refresh container rows on input
	stream := i.Stream()
	go func() {
		for s := range stream {
			filter = s
			RebuildRows()
			ui.Render(i)
		}
	}()

	i.InputHandlers()
	ui.Handle("/sys/kbd/<escape>", func(ui.Event) {
		filter = ""
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		filter = i.Data
		ui.StopLoop()
	})
	ui.Loop()
	RebuildRows()
}
