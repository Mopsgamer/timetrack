package timet

import (
	"errors"
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
	RecordActionNormal  = "normal"
	RecordActionAdded   = "added"
	RecordActionDeleted = "deleted"
)

var DefaultData = Data{
	Records: []*Record{},
}

const DateFormatReadable string = time.DateTime
const DateFormat string = time.RFC3339

func MakeRecord(name string, time time.Time) *Record {
	if name == "" {
		return nil
	}
	return &Record{Name: name, Since: time.Format(DateFormat)}
}

func MakeRecordActed(record *Record) *RecordActed {
	if record == nil {
		return nil
	}
	return &RecordActed{Record: record, Action: RecordActionNormal}
}

func MakeRecordActedList(records []*Record) []*RecordActed {
	recordsFm := make([]*RecordActed, len(records))
	for reci, rec := range records {
		recordsFm[reci] = MakeRecordActed(rec)
	}
	return recordsFm
}

func MakeRecordFormat(recordAc *RecordActed, recordAcIndex int) *RecordFormat {
	if recordAc == nil {
		return nil
	}
	recordDate, errRecordDate := time.Parse(DateFormat, recordAc.Since)
	if errRecordDate != nil {
		return nil
	}
	colorize := MakeColorizer(recordAc.Action)
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
	}
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

func MakeColorizer(action string) func(string) string {
	color := ansi.ColorFunc(map[string]string{
		RecordActionNormal:  "white",
		RecordActionAdded:   "green+h",
		RecordActionDeleted: "red+h",
	}[action])
	return color
}

func RecordsActedToRows(recordsAc []*RecordActed) ([][]string, error) {
	rows := [][]string{}
	rmCount := 0
	if recordsAc == nil {
		return nil, errors.New("bad list")
	}
	for recordAcIndex, recordAc := range recordsAc {
		if recordAc == nil {
			return nil, errors.New("bad record")
		}
		isDeleted := recordAc.Action == RecordActionDeleted
		i := recordAcIndex + 1 - rmCount
		if isDeleted {
			i = -1
			rmCount++
		}
		recordFm := MakeRecordFormat(recordAc, i)
		if recordFm == nil {
			return nil, errors.New("bad record, can not format")
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
