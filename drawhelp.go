package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

func splitTextIntoLines(text string, limitX int) []string {
	result := []string{}
	slice := ""
	current := 0
	chars := strings.SplitSeq(text, "")
	for r := range chars {
		current += 1
		if r == "\n" || current >= limitX {
			result = append(result, slice)
			slice = ""
			current = 0
		}
		if r != "\n" {
			slice += r
		}
	}
	result = append(result, slice)
	return result
}

func drawParagraph(screen tcell.Screen, size tcell.WindowSize, text string, style tcell.Style) {
	lines := splitTextIntoLines(text, size.PixelWidth)
	lines = lines[:min(size.PixelHeight, len(lines))]
	for i, line := range lines {
		drawText(screen, size.Width, size.Height+i, line, style)
	}
}
