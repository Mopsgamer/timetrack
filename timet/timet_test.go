package timet

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mgutz/ansi"
	"github.com/stretchr/testify/assert"
)

func TestTimet(t *testing.T) {
	check := assert.New(t)
	t.Run("MakeRecord", func(t *testing.T) {
		now := time.Now()
		record := MakeRecord("", now)
		check.Nil(record)
		record = MakeRecord("test", now)
		check.NotNil(record)
	})
	t.Run("MakeRecordActed", func(t *testing.T) {
		now := time.Now()
		recordAc := MakeRecordActed(MakeRecord("", now))
		check.Nil(recordAc)
		recordAc = MakeRecordActed(MakeRecord("test", now))
		check.NotNil(recordAc)
	})
	t.Run("MakeRecordActedList", func(t *testing.T) {
		now := time.Now()
		MakeRecordActedList([]*Record{
			MakeRecord("", now),
		})
		MakeRecordActedList([]*Record{
			MakeRecord("test", now),
		})
		MakeRecordActedList([]*Record{
			MakeRecord("", now),
			MakeRecord("test", now),
		})
	})
	t.Run("MakeRecordFormat", func(t *testing.T) {
		now := time.Now()
		recordFm := MakeRecordFormat(MakeRecordActed(MakeRecord("", now)), -1)
		check.Nil(recordFm)
		recordFm = MakeRecordFormat(MakeRecordActed(MakeRecord("test", now)), -1)
		check.NotNil(recordFm)
	})
	t.Run("MakeColorizer", func(t *testing.T) {
		text := "abc"
		check.Equal(MakeColorizer(RecordActionNormal)(text), ansi.ColorFunc("white")(text))
		check.Equal(MakeColorizer(RecordActionAdded)(text), ansi.ColorFunc("green+h")(text))
		check.Equal(MakeColorizer(RecordActionDeleted)(text), ansi.ColorFunc("red+h")(text))
	})
	t.Run("ParseFile", func(t *testing.T) {
		dataTestFile := filepath.Join(PathRoot, "parsefile.test.json")
		defer os.Remove(dataTestFile)

		testData := []struct {
			data      string
			wantError bool
		}{
			{
				data:      `{"Records":[{"Name":"test timestamp","Since":"` + time.Now().Format(DateFormat) + `"}]}`,
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
				_, err := ParseFile(dataTestFile)
				check.Error(err)
				continue
			}
			_, err := ParseFile(dataTestFile)
			check.Nil(err)
		}
	})
}
