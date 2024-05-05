package texttable

import "strings"

const (
	AlignRight = iota
	AlignLeft  = iota
)

type TextTableOptions struct {
	AlignH       []int
	StringLength func(string) int
}

func Make(rows [][]string, options TextTableOptions) string {

	stringLength := options.StringLength
	alignListH := options.AlignH

	result := ""
	colSizes := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, col := range row {
			if stringLength(col) > colSizes[i] {
				colSizes[i] = stringLength(col)
			}
		}
	}
	for rowi, row := range rows {
		for coli, col := range row {
			colLen := stringLength(col)
			alignH := alignListH[coli]
			if alignH == AlignLeft {
				padSize := colSizes[coli] - colLen
				if coli < len(row)-1 {
					padSize += 2
				}
				pad := strings.Repeat(" ", padSize)
				result += col + pad
			} else if alignH == AlignRight {
				padSize := colSizes[coli] - colLen
				if coli > 0 {
					padSize += 2
				}
				pad := strings.Repeat(" ", padSize)
				result += pad + col
			}
		}
		if rowi != len(rows) {
			result += "\n"
		}
	}
	return result
}
