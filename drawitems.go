package main

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func drawItems(screen tcell.Screen, items []*Item, current int, highlight bool, size tcell.WindowSize) {
	half := size.PixelHeight / 2
	from, to := 0, len(items)
	if current > half {
		if current > len(items)-half {
			from, to = len(items)-size.PixelHeight, len(items)
		} else {
			from, to = current-half-1, current+half
		}
	}

	for i, item := range items[from:to] {
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		since := FormatSince(time.Since(item.Since))
		if highlight && item == items[current] {
			style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGhostWhite)
			drawText(screen, size.Width+len(item.Name), i+size.Height, strings.Repeat(" ", max(0, size.PixelWidth-size.Width-len(item.Name)-len(since))), style)
		}
		drawText(screen, size.Width, i+size.Height, item.Name, style)
		drawTextRight(screen, size.PixelWidth-size.Width+1, i+size.Height, since, style)
	}
}
