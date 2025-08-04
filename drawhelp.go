package main

import (
	"github.com/gdamore/tcell/v2"
)

func splitTextIntoLines(text string, limit int) []string {
	result := []string{}
	slice := ""
	current := 0
	for _, r := range text {
		current += 1
		if r == '\n' || current >= limit {
			result = append(result, slice)
			slice = ""
			current = 0
			continue
		}
		slice += string(r)
	}
	return result
}

func drawParagraph(screen tcell.Screen, size tcell.WindowSize, text string, style tcell.Style) {
	lines := splitTextIntoLines(text, size.PixelWidth)

	for i, line := range lines[:min(size.PixelHeight, len(lines)-1)] {
		drawText(screen, size.Width, size.Height+i, line, style)
	}
}
