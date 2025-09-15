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

	i := 0
	for _, item := range items[from:to] {
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		since := FormatSince(time.Since(item.Since))
		isCurrent := item == items[current]
		if highlight && isCurrent {
			style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGhostWhite)
			drawText(screen, size.Width+len(item.Name), i+size.Height, strings.Repeat(" ", max(0, size.PixelWidth-size.Width-len(item.Name)-len(since))), style)
		}

		drawText(screen, size.Width, i+size.Height, item.Name, style)
		drawTextRight(screen, size.PixelWidth-size.Width+1, i+size.Height, since, style)
		if isCurrent {
			date := item.Since.Format("2006-01-02 15:04:05")
			i += 1
			drawText(screen, size.Width, i+size.Height, strings.Repeat(" ", max(0, size.PixelWidth-size.Width-len(date))), style)
			drawTextRight(screen, size.PixelWidth-size.Width+1, i+size.Height, date, style)
		}
		i += 1
	}
}
