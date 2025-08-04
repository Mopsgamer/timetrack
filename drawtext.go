package main

import (
	"github.com/gdamore/tcell/v2"
)

func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, r := range text {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

func drawTextRight(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	drawText(screen, x-len(text)+1, y, text, style)
}
