package main

import (
	"github.com/gdamore/tcell/v2"
)

func redraw(screen tcell.Screen, state State) {
	screen.Clear()
	w, h := screen.Size()
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	listBox := tcell.WindowSize{PixelWidth: w - 2, PixelHeight: h - 2, Height: 1, Width: 2}
	if state.Window == StateSearch || state.Window == StateNew {
		listBox.PixelHeight = h - 3
	}

	switch state.Window {
	case StateSearch:
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
		x := 0
		y := h - 1
		label := "New: "
		drawText(screen, x, y, label+state.NewItem.Name, style)
		x = len(label) + state.InputNew.Cursor
		screen.ShowCursor(x, y)
	case StateHelp:
		drawParagraph(screen, listBox, help, style)
		screen.HideCursor()
	case StateList:
		screen.HideCursor()
	}

	if len(state.ItemsFound) == 0 {
		drawParagraph(screen, listBox, help, style)
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
