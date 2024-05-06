package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
	"timet/timet"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	check := assert.New(t)
	t.Run("Data", func(t *testing.T) {
		t.Run("\"\" Manager", func(t *testing.T) {
			m := New("")
			errLoad := m.DataLoadFromFile()
			check.Error(errLoad)
			errSave := m.DataSaveToFile()
			check.Error(errSave)
		})
		t.Run("\"data_abcd.test.json\" Manager", func(t *testing.T) {
			dataTestFile := filepath.Join(timet.PathRoot, "data_abcd.test.json")
			defer os.Remove(dataTestFile)
			m := New(dataTestFile)
			errLoad := m.DataLoadFromFile()
			check.Error(errLoad)
			errSave := m.DataSaveToFile()
			check.Nil(errSave)
		})
		t.Run("\"data.test.json\" Manager", func(t *testing.T) {
			dataTestFile := filepath.Join(timet.PathRoot, "data.test.json")
			defer os.Remove(dataTestFile)
			m := New(dataTestFile)

			testData := []struct {
				data      string
				wantError bool
			}{
				{
					data:      `{"Records":[{"Name":"test timestamp","Since":"` + time.Now().Format(timet.DateFormat) + `"}]}`,
					wantError: false,
				},
				{
					data:      `[]`,
					wantError: true,
				},
				{
					data:      `{}`,
					wantError: false,
				},
				{
					data:      `{`,
					wantError: true,
				},
				{
					data:      ``,
					wantError: true,
				},
			}

			for _, data := range testData {
				os.WriteFile(dataTestFile, []byte(data.data), 0644)
				if data.wantError {
					check.Error(m.DataLoadFromFile())
					continue
				}
				check.Nil(m.DataLoadFromFile())
			}

			errSave := m.DataSaveToFile()
			check.Nil(errSave)
		})
	})
	t.Run("Actions", func(t *testing.T) {
		t.Run("Limit", func(t *testing.T) {
			m := New("")
			for i := range timet.MaxRecords {
				m.Create(fmt.Sprint(i), "", false)
			}
			check.Equal(len(m.Data.Records), timet.MaxRecords)
			_, errGotLimit := m.Create("?", "", false)
			check.Error(errGotLimit)
		})
		t.Run("List", func(t *testing.T) {
			m := New("")
			_, errListEmpty := m.List("")
			check.Nil(errListEmpty)
			m.Create("test", "", false)
			_, errListAll := m.List("")
			check.Nil(errListAll)
			_, errListSingle := m.List("test")
			check.Nil(errListSingle)
		})
		t.Run("Create", func(t *testing.T) {
			m := New("")
			_, errName := m.Create("", "", false)
			check.Error(errName)
			_, errNoDate := m.Create("no date above", "", false)
			check.Nil(errNoDate)
			_, errBadDate := m.Create("bad date", "dsfsd", false)
			check.Error(errBadDate)
			_, errGoodDate := m.Create("good date", "2004-12-25 23:00:00", false)
			check.Nil(errGoodDate)
			_, errBelow := m.Create("below", "", true)
			check.Nil(errBelow)
			check.Equal(m.Data.Records[len(m.Data.Records)-1].Name, "below")
			_, errAbove := m.Create("above", "", false)
			check.Nil(errAbove)
			check.Equal(m.Data.Records[0].Name, "above")
		})
		t.Run("Remove", func(t *testing.T) {
			m := New("")
			_, errRmEmpty := m.Remove("")
			check.Nil(errRmEmpty)
			m.Create("test", "", false)
			_, errRmAll := m.Remove("")
			check.Nil(errRmAll)
			check.Equal(len(m.Data.Records), 0)
		})
		t.Run("Reset", func(t *testing.T) {
			m := New("")
			_, err := m.Reset()
			check.Nil(err)
		})
	})
	t.Run("Finders", func(t *testing.T) {
		t.Run("FindRecordAll", func(t *testing.T) {
			m := New("")
			check.Equal(len(m.FindRecordAll(nil, "")), 0)
			m.Create("hello", "", true)
			m.Create("hell", "", true)
			check.Equal(len(m.FindRecordAll(m.Data.Records, "hello")), 1)
			check.Equal(len(m.FindRecordAll(m.Data.Records, "h?l?o")), 1)
			check.Equal(len(m.FindRecordAll(m.Data.Records, "*o")), 1)
			check.Equal(len(m.FindRecordAll(m.Data.Records, "h*")), 2)
			check.Equal(len(m.FindRecordAll(m.Data.Records, "h*o")), 1)
			check.Equal(len(m.FindRecordAll(m.Data.Records, "h**o")), 1)
		})
		t.Run("FindRecordIndex", func(t *testing.T) {
			m := New("")
			check.Equal(m.FindRecordIndex(nil, ""), -1)
			m.Create("hello", "", true)
			check.Equal(m.FindRecordIndex(m.Data.Records, "hello"), 0)
			check.Equal(m.FindRecordIndex(m.Data.Records, "h?l?o"), 0)
			check.Equal(m.FindRecordIndex(m.Data.Records, "*o"), 0)
			check.Equal(m.FindRecordIndex(m.Data.Records, "h*"), 0)
			check.Equal(m.FindRecordIndex(m.Data.Records, "h*o"), 0)
			check.Equal(m.FindRecordIndex(m.Data.Records, "h**o"), 0)
		})
		t.Run("FindRecord", func(t *testing.T) {
			m := New("")
			check.Nil(m.FindRecord(nil, ""))
			m.Create("hello", "", true)
			check.NotNil(m.FindRecord(m.Data.Records, "hello"))
			check.NotNil(m.FindRecord(m.Data.Records, "h?l?o"))
			check.NotNil(m.FindRecord(m.Data.Records, "*o"))
			check.NotNil(m.FindRecord(m.Data.Records, "h*"))
			check.NotNil(m.FindRecord(m.Data.Records, "h*o"))
			check.NotNil(m.FindRecord(m.Data.Records, "h**o"))
		})
	})
}
