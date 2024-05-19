package timet

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mopsgamer/timetrack/texttable"

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
	Records []Record
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

type IRowFormatable interface {
	GetAction() string
	Format(index int) *RecordFormat
}

const (
	RecordActionNormal  = "normal"
	RecordActionAdded   = "added"
	RecordActionDeleted = "deleted"
)

var DefaultData = Data{
	Records: []Record{},
}

const DateFormatReadable string = time.DateTime
const DateFormat string = time.RFC3339

func MakeRecord(name string, time time.Time) *Record {
	if name == "" {
		return nil
	}
	return &Record{Name: name, Since: time.Format(DateFormat)}
}

func MakeRecordActedList(records []Record) []RecordActed {
	recordsFm := make([]RecordActed, len(records))
	for reci, rec := range records {
		recordsFm[reci] = RecordActed{Record: &rec, Action: rec.GetAction()}
	}
	return recordsFm
}

func (record *RecordActed) GetAction() string {
	return record.Action
}

func (record *Record) GetAction() string {
	return RecordActionNormal
}

func (record *Record) Format(recordIndex int) *RecordFormat {
	if record == nil {
		return nil
	}
	recordDate, errRecordDate := time.Parse(DateFormat, record.Since)
	if errRecordDate != nil {
		return nil
	}
	colorize := MakeColorizer(record.GetAction())
	var index string
	if recordIndex < 0 {
		index = "x"
	} else {
		index = fmt.Sprint(recordIndex)
	}
	index = colorize(index)
	name := colorize(record.Name)
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

func StringRows(formatableList []IRowFormatable) ([][]string, error) {
	rows := [][]string{}
	rmCount := 0
	if formatableList == nil {
		return nil, errors.New("bad list")
	}
	for formatableIndex, formatable := range formatableList {
		if formatable == nil {
			return nil, errors.New("bad record")
		}
		isDeleted := formatable.GetAction() == RecordActionDeleted
		i := formatableIndex + 1 - rmCount
		if isDeleted {
			i = -1
			rmCount++
		}
		recordFm := formatable.Format(i)
		rows = append(rows, []string{
			recordFm.Index,
			recordFm.Name,
			recordFm.Since,
			recordFm.Date,
		})
	}
	return rows, nil
}

func String(rowFormatableList []IRowFormatable) (string, error) {
	head := []string{"#", "Name", "Since", "Date"}
	for rowi, row := range head {
		head[rowi] = ansi.Color(row, "white+u")
	}

	rows, err := StringRows(rowFormatableList)
	if err != nil {
		return "", err
	}

	rows = append([][]string{head}, rows...)

	return texttable.Make(rows, textTableOptions), nil
}
