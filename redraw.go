package main

import (
	"github.com/gdamore/tcell/v2"
)

var listBox = tcell.WindowSize{}

func redraw(screen tcell.Screen, state *State) {
	screen.Clear()
	w, h := screen.Size()
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	listBox = tcell.WindowSize{PixelWidth: w - 2, PixelHeight: h - 2, Height: 1, Width: 2}

	switch state.Window {
	case StateSearch:
		listBox.PixelHeight = h - 3
		x := 0
		y := h - 1
		label := "Search: /"
		drawText(screen, x, y, label, style)
		x = len(label) - 1
		if _, err := state.SearchRegexp(); err != nil {
			style = style.Foreground(tcell.ColorDarkRed)
		}
		drawText(screen, x, y, "/"+state.SearchContent+"/", style)
		x = len(label) + state.InputSearch.Cursor
		screen.ShowCursor(x, y)
	case StateNew:
		listBox.PixelHeight = h - 3
		x := 0
		y := h - 1
		label := "New: "
		drawText(screen, x, y, label+state.InputNew.Value, style)
		x = len(label) + state.InputNew.Cursor
		screen.ShowCursor(x, y)
	case StateRename:
		listBox.PixelHeight = h - 3
		x := 0
		y := h - 1
		label := "Rename: "
		drawText(screen, x, y, label+state.InputRename.Value, style)
		x = len(label) + state.InputRename.Cursor
		screen.ShowCursor(x, y)
	case StateHelp:
		drawParagraph(screen, listBox, help, state.HelpScroll, style)
		screen.HideCursor()
	case StateList:
		screen.HideCursor()
	}

	if len(state.ItemsFound) == 0 {
		drawParagraph(screen, listBox, help, state.HelpScroll, style)
	} else if state.Window != StateHelp {
		drawItems(screen, state.ItemsFound, state.ItemFound, true, listBox)
	}

	for x := range w {
		screen.SetContent(x, listBox.Height-1, '─', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		screen.SetContent(x, listBox.PixelHeight+1, '━', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
	for y := listBox.Height; y < h-1; y++ {
		screen.SetContent(0, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		screen.SetContent(w-1, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
	screen.SetContent(0, listBox.Height-1, '┌', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	screen.SetContent(w-1, listBox.Height-1, '┐', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	screen.SetContent(0, listBox.PixelHeight+1, '┕', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	screen.SetContent(w-1, listBox.PixelHeight+1, '┙', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))

	screen.Show()
}
