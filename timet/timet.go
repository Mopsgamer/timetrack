package timet

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"timet/prettyms"
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
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return DefaultData, err
	}
	return parsed, nil
}

func String(records []*RecordActed) (string, error) {
	// Clone records
	cloned := make([]*RecordActed, len(records))
	copy(cloned, records)

	// Header row
	head := []string{"#", "Name", "Since", "Date"}
	rows := [][]string{head}
	for rowi, row := range head {
		head[rowi] = ansi.Color(row, "white+u")
	}

	// Format switch for different actions
	formatSwitch := map[string]func(string) string{
		RowActionNormal:  func(s string) string { return s },
		RowActionAdded:   func(s string) string { return ansi.Color(s, "green+h") }, // green color
		RowActionDeleted: func(s string) string { return ansi.Color(s, "red+h") },   // red color
	}

	// Iterate over the records
	for recordIndex := 0; recordIndex < len(cloned); recordIndex++ {
		record := cloned[recordIndex]
		if record.Since == "" {
			record.Since = time.Now().Format(DateFormat)
		}
		recordDate, errRecordDate := time.Parse(DateFormat, record.Since)
		if errRecordDate != nil {
			return "", errors.New(`invalid date "` + record.Since + `"`)
		}
		isRemoved := record.Action == RowActionDeleted

		// Color by action
		var colorize func(string) string
		if format, ok := formatSwitch[record.Action]; ok {
			colorize = format
		} else {
			colorize = formatSwitch[RowActionNormal]
		}

		// Add row
		var rowN = colorize(fmt.Sprint(recordIndex + 1))
		if isRemoved {
			rowN = colorize("x")
		}
		rowSinceTime, err := time.Parse(DateFormat, record.Since)
		if err != nil {
			return "", err
		}
		var recordName = colorize(record.Name)
		var rowSince = colorize(prettyms.Time(rowSinceTime))
		var rowDate = colorize(recordDate.Format(DateFormatReadable))
		var row = []string{rowN, recordName, rowSince, rowDate}
		rows = append(rows, row)

		// If removed, remove from cloned and decrement index
		if isRemoved {
			cloned = append(cloned[:recordIndex], cloned[recordIndex+1:]...)
			recordIndex--
		}
	}

	return texttable.Make(rows, textTableOptions), nil
}
