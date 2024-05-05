package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
	"timet/timet"

	"github.com/gobwas/glob"
)

var (
	messageEmpty       = `no records`
	messageNoMatches   = `no matches`
	messageLimit       = `you have reached the limit of ` + fmt.Sprint(timet.MaxRecords) + ` records`
	messageInvalidName = `invalid name`
	// messageInvalidName2       = `invalid second name`
	messageExistsName = `record with this name already exists`
	// messageSuccSwappedNames   = `successfully swapped names`
	// messageSuccSwappedPlaces  = `successfully swapped places`
	// messageSuccMoved          = func(direction string, count string) string { return `moved ` + direction + ` by ${count}` }
	// messageCantMove           = func(direction string) string { return `can not move` + direction }
	// messageInvalidInteger     = func(varname string) string { return `'` + varname + `' must be an integer` }
	// messageInvalidGreaterLess = func(varname string, greater string, less string) string {
	// return `'${varname}' must be greater than ${greater} and less than ${less}`
	// }
)

type Manager struct {
	Path string
	Data timet.Data
}

func New(path string) *Manager {
	return &Manager{Path: path}
}

func (m *Manager) DataLoadFromFile() error {
	data, err := timet.ParseFile(m.Path)
	if err == nil {
		m.Data = data
	}
	return err
}

func (m *Manager) DataJsonMarshal() ([]byte, error) {
	str, err := json.Marshal(m.Data)
	if err != nil {
		return nil, err
	}
	return str, nil
}

func (m *Manager) DataSaveToFile() error {
	if m.Path == "" {
		return errors.New("Manager is virtual")
	}
	str, err := m.DataJsonMarshal()
	if err != nil {
		return err
	}
	err = os.WriteFile(m.Path, str, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) FindRecord(records []*timet.Record, name string) *timet.Record {
	recordIndex := m.FindRecordIndex(records, name)
	if recordIndex == -1 {
		return nil
	}
	record := records[recordIndex]
	return record
}

func (m *Manager) FindRecordIndex(records []*timet.Record, name string) int {
	for i, v := range records {
		if glob.MustCompile(name).Match(v.Name) {
			return i
		}
	}
	return -1
}

func (m *Manager) FindRecordAll(records []*timet.Record, name string) []*timet.Record {
	result := []*timet.Record{}
	for _, v := range records {
		if glob.MustCompile(name).Match(v.Name) {
			result = append(result, v)
		}
	}
	return result
}

func (m *Manager) List(name string) (string, error) {
	if len(m.Data.Records) == 0 {
		return messageEmpty, nil
	}
	if name == "" {
		name = "**"
	}
	recordsF := m.FindRecordAll(m.Data.Records, name)
	if len(recordsF) == 0 {
		return "", errors.New(messageNoMatches)
	}
	return timet.String(timet.FormatRecordList(m.Data.Records))
}

func (m *Manager) Create(name string, date string, below bool) (string, error) {
	if len(m.Data.Records) >= timet.MaxRecords {
		return "", errors.New(messageLimit)
	}
	if name == "" {
		return "", errors.New(messageInvalidName)
	}
	if date == "" {
		date = time.Now().Format(timet.DateFormatReadable)
	}
	recordTime, err := time.ParseInLocation(timet.DateFormatReadable, date, time.Local)
	if err != nil {
		return "", err
	}
	similar := m.FindRecord(m.Data.Records, name)
	if similar != nil {
		return "", errors.New(messageExistsName)
	}
	record := timet.CreateRecord(name, recordTime)
	if below {
		m.Data.Records = append(m.Data.Records, &record)
	} else {
		m.Data.Records = append([]*timet.Record{&record}, m.Data.Records...)
	}
	m.DataSaveToFile()
	recordsFm := timet.FormatRecordList(m.Data.Records)
	var recordFm *timet.RecordActed
	if below {
		recordFm = recordsFm[0]
	} else {
		recordFm = recordsFm[len(recordsFm)-1]
	}
	recordFm.Action = timet.RowActionAdded
	return timet.String(recordsFm)
}

func (m *Manager) Remove(name string) (string, error) {
	if len(m.Data.Records) == 0 {
		return messageEmpty, nil
	}
	if name == "" {
		name = "**"
	}
	similar := m.FindRecord(m.Data.Records, name)
	if similar == nil {
		return "", errors.New(messageNoMatches)
	}
	recordsFm := timet.FormatRecordList(m.Data.Records)
	var counterRm = 0 // needed for right painting
	for {
		recordIndexRm := m.FindRecordIndex(m.Data.Records, name)
		if recordIndexRm < 0 {
			break
		}
		recordsFm[recordIndexRm+counterRm].Action = timet.RowActionDeleted
		// remove el by index (recordIndexRm)
		m.Data.Records = append(m.Data.Records[:recordIndexRm], m.Data.Records[recordIndexRm+1:]...)
		counterRm++
	}
	m.DataSaveToFile()
	return timet.String(recordsFm)
}

func (m *Manager) Reset() (string, error) {
	m.Data = timet.DefaultData
	if m.Path != "" {
		err := m.DataSaveToFile()
		if err != nil {
			return "", err
		}
	}
	return "reset completed", nil
}
