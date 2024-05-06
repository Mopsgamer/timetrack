package texttable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextTable(t *testing.T) {
	check := assert.New(t)
	t1Row := Make([][]string{{"col1", "col2"}}, TextTableOptions{})
	check.Equal("col1  col2", t1Row)
	t2Rows := Make([][]string{{"11", "12"}, {"21 hello", "22 h"}}, TextTableOptions{})
	check.Equal("      11    12\n21 hello  22 h", t2Rows)
	tWithParams := Make([][]string{{"11", "12"}, {"21 hello", "22 h"}}, TextTableOptions{AlignH: []int{AlignLeft, AlignRight}})
	check.Equal("11          12\n21 hello  22 h", tWithParams)
	tPartial := Make([][]string{{"11"}, {"21 hello", "22 h"}}, TextTableOptions{})
	check.Equal("      11\n21 hello  22 h", tPartial)
}

func BenchmarkTextTable(b *testing.B) {
	rows := [][]string{}
	for n := 0; n < b.N; n++ {
		rows = append(rows, []string{"11", "12", "13"})
	}
	Make(rows, TextTableOptions{})
}
