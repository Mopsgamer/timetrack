package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/mopsgamer/timetrack/timet"

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

type ManagerObserver struct {
	onCreate func(record timet.Record, recordIndex int)
	onRemove func(record timet.Record, recordIndex int)
	onReset  func()
}

func (ob *ManagerObserver) Follow(m *Manager) {
	if ob.IsFollowed(m) {
		return
	}
	m.observers = append(m.observers, *ob)
}

func (ob *ManagerObserver) Unfollow(m *Manager) {
	obListLen := len(m.observers)
	for i, observer := range m.observers {
		if ob == &observer {
			m.observers[obListLen-1], m.observers[i] = m.observers[i], m.observers[obListLen-1]
			m.observers = m.observers[:obListLen-1]
			return
		}
	}
}

func (ob *ManagerObserver) IsFollowed(m *Manager) bool {
	for _, observer := range m.observers {
		if ob == &observer {
			return true
		}
	}
	return false
}

type Manager struct {
	observers []ManagerObserver
	Path      string
	Data      timet.Data
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
	errWrite := os.WriteFile(m.Path, str, 0644)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

func FindRecord(records []timet.Record, name string) *timet.Record {
	recordIndex := FindRecordIndex(records, name)
	if recordIndex == -1 {
		return nil
	}
	record := records[recordIndex]
	return &record
}

func FindRecordIndex(records []timet.Record, name string) int {
	for i, v := range records {
		if glob.MustCompile(name).Match(v.Name) {
			return i
		}
	}
	return -1
}

func FindRecordAll(records []timet.Record, name string) []timet.Record {
	result := []timet.Record{}
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
	recordsF := FindRecordAll(m.Data.Records, name)
	if len(recordsF) == 0 {
		return "", errors.New(messageNoMatches)
	}
	rows := make([]timet.IRowFormatable, len(recordsF))
	for i, v := range recordsF {
		rows[i] = &v
	}
	return timet.String(rows)
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
	similar := FindRecord(m.Data.Records, name)
	if similar != nil {
		return "", errors.New(messageExistsName)
	}
	record := *timet.MakeRecord(name, recordTime)
	recordIndex := 0
	if below {
		m.Data.Records = append(m.Data.Records, record)
	} else {
		recordIndex = len(m.Data.Records)
		m.Data.Records = append([]timet.Record{record}, m.Data.Records...)
	}
	for _, ob := range m.observers {
		ob.onCreate(record, recordIndex)
	}
	defer m.DataSaveToFile()
	recordsFm := timet.MakeRecordActedList(m.Data.Records)
	var recordFm timet.RecordActed
	if below {
		recordFm = recordsFm[len(recordsFm)-1]
	} else {
		recordFm = recordsFm[0]
	}
	recordFm.Action = timet.RecordActionAdded
	rows := make([]timet.IRowFormatable, len(recordsFm))
	for i, v := range recordsFm {
		rows[i] = &v
	}
	return timet.String(rows)
}

func (m *Manager) Remove(name string) (string, error) {
	if len(m.Data.Records) == 0 {
		return messageEmpty, nil
	}
	if name == "" {
		name = "**"
	}
	similar := FindRecord(m.Data.Records, name)
	if similar == nil {
		return "", errors.New(messageNoMatches)
	}
	recordsFm := timet.MakeRecordActedList(m.Data.Records)
	var counterRm = 0 // needed for right painting
	for {
		recordIndexRm := FindRecordIndex(m.Data.Records, name)
		if recordIndexRm < 0 {
			break
		}
		recordRm := m.Data.Records[recordIndexRm]
		recordsFm[recordIndexRm+counterRm].Action = timet.RecordActionDeleted
		// remove el by index (recordIndexRm)
		m.Data.Records = append(m.Data.Records[:recordIndexRm], m.Data.Records[recordIndexRm+1:]...)
		counterRm++
		for _, ob := range m.observers {
			ob.onRemove(recordRm, recordIndexRm)
		}
	}
	defer m.DataSaveToFile()
	rows := make([]timet.IRowFormatable, len(recordsFm))
	for i, v := range recordsFm {
		rows[i] = &v
	}
	return timet.String(rows)
}

func (m *Manager) Reset() (string, error) {
	m.Data = timet.DefaultData
	for _, ob := range m.observers {
		ob.onReset()
	}
	defer m.DataSaveToFile()
	return "reset completed", nil
}
