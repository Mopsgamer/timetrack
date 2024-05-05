package timet

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"timet/texttable"

	"encoding/json"

	"github.com/acarl005/stripansi"
	"github.com/mgutz/ansi"
)

var (
	PathRoot   string = filepath.Dir(".")
	PathData   string = filepath.Join(PathRoot, "data.json")
	Version    string = "v0.0.1"
	MaxRecords int    = 200
)

var textTableOptions = texttable.TextTableOptions{
	AlignH: []int{
		texttable.AlignRight,
		texttable.AlignRight,
		texttable.AlignRight,
		texttable.AlignRight,
	},
	StringLength: func(s string) int {
		return len(stripansi.Strip(s))
	},
}

type Data struct {
	Records []*Record
}

type Record struct {
	Name  string
	Since string
}

type RecordActed struct {
	*Record
	Action string
}

type RecordFormat struct {
	Index string
	Name  string
	Since string
	Date  string
}

const (
	RowActionNormal  = "normal"
	RowActionAdded   = "added"
	RowActionDeleted = "deleted"
)

var DefaultData = Data{
	Records: []*Record{},
}

const DateFormatReadable string = time.DateTime
const DateFormat string = time.RFC3339

func CreateRecord(name string, time time.Time) Record {
	return Record{Name: name, Since: time.Format(DateFormat)}
}

func FormatRecordList(records []*Record) []*RecordActed {
	recordsFm := make([]*RecordActed, len(records))
	for reci, rec := range records {
		recordsFm[reci] = &RecordActed{Record: rec, Action: RowActionNormal}
	}
	return recordsFm
}

func ParseFile(path string) (Data, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultData, err
	}
	var parsed Data
	errParse := json.Unmarshal(data, &parsed)
	if errParse != nil {
		return DefaultData, errParse
	}
	return parsed, nil
}

func NewColorizer(action string) func(string) string {
	color := (map[string]func(string) string{
		RowActionNormal:  func(s string) string { return s },
		RowActionAdded:   func(s string) string { return ansi.Color(s, "green+h") },
		RowActionDeleted: func(s string) string { return ansi.Color(s, "red+h") },
	}[action])
	return color
}

func (recordAc *RecordActed) Format(recordAcIndex int) (*RecordFormat, error) {
	recordDate, errRecordDate := time.Parse(DateFormat, recordAc.Since)
	if errRecordDate != nil {
		return nil, errRecordDate
	}
	colorize := NewColorizer(recordAc.Action)
	var index string
	if recordAcIndex < 0 {
		index = "x"
	} else {
		index = fmt.Sprint(recordAcIndex)
	}
	index = colorize(index)
	name := colorize(recordAc.Name)
	since := colorize(fmt.Sprintf("%v", time.Since(recordDate)))
	date := colorize(recordDate.Format(DateFormatReadable))

	return &RecordFormat{
		Index: index,
		Name:  name,
		Since: since,
		Date:  date,
	}, nil
}

func RecordsActedToRows(recordsAc []*RecordActed) ([][]string, error) {
	rows := [][]string{}
	rmCount := 0
	for recordAcIndex, recordAc := range recordsAc {
		isDeleted := recordAc.Action == RowActionDeleted
		i := recordAcIndex + 1 - rmCount
		if isDeleted {
			i = -1
			rmCount++
		}
		recordFm, err := recordAc.Format(i)
		if err != nil {
			return rows, err
		}
		rows = append(rows, []string{
			recordFm.Index,
			recordFm.Name,
			recordFm.Since,
			recordFm.Date,
		})
	}
	return rows, nil
}

func String(recordsAc []*RecordActed) (string, error) {
	cloned := make([]*RecordActed, len(recordsAc))
	copy(cloned, recordsAc)

	head := []string{"#", "Name", "Since", "Date"}
	for rowi, row := range head {
		head[rowi] = ansi.Color(row, "white+u")
	}

	rows, err := RecordsActedToRows(cloned)
	if err != nil {
		return "", err
	}

	rows = append([][]string{head}, rows...)

	return texttable.Make(rows, textTableOptions), nil
}
