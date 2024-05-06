package texttable

import (
	"strings"

	"github.com/acarl005/stripansi"
)

const (
	AlignRight = iota
	AlignLeft  = iota
)

type TextTableOptions struct {
	AlignH       []int
	StringLength func(string) int
}

func ColumnFit(cell string, fitSize int, stringLength func(string) int, alignH int) string {
	colLen := stringLength(cell)
	padSize := fitSize - colLen
	pad := strings.Repeat(" ", padSize)
	if alignH == AlignLeft {
		return cell + pad
	}
	return pad + cell
}

func ColumnMaxSizes(rows [][]string, stringLength func(string) int) []int {
	colSizes := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, col := range row {
			if i >= len(colSizes) {
				colSizes = append(colSizes, 0)
			}
			if stringLength(col) > colSizes[i] {
				colSizes[i] = stringLength(col)
			}
		}
	}
	return colSizes
}

func Make(rows [][]string, options TextTableOptions) string {
	stringLength := options.StringLength
	if stringLength == nil {
		stringLength = func(s string) int { return len(stripansi.Strip(s)) }
	}
	alignListH := options.AlignH
	result := ""
	colSizes := ColumnMaxSizes(rows, stringLength)
	for rowi, row := range rows {
		for coli, cell := range row {
			var alignH int = AlignRight
			if alignListH != nil && coli < len(alignListH) {
				alignH = alignListH[coli]
			}
			result += ColumnFit(cell, colSizes[coli], stringLength, alignH)
			if coli < len(row)-1 {
				result += "  "
			}
		}
		if rowi < len(rows)-1 {
			result += "\n"
		}
	}
	return result
}
