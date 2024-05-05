package main

import (
	"fmt"
	"testing"
	"timet/manager"
	"timet/timet"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	check := assert.New(t)
	t.Run("Limit", func(t *testing.T) {
		m := manager.New("")
		for i := range timet.MaxRecords {
			m.Create(fmt.Sprint(i), "", false)
		}
		check.Equal(len(m.Data.Records), timet.MaxRecords)
		_, errGotLimit := m.Create("?", "", false)
		check.Error(errGotLimit)
	})
	t.Run("Actions", func(t *testing.T) {
		t.Run("List", func(t *testing.T) {
			m := manager.New("")
			_, errListEmpty := m.List("")
			check.Nil(errListEmpty)
			m.Create("test", "", false)
			_, errListAll := m.List("")
			check.Nil(errListAll)
			_, errListSingle := m.List("test")
			check.Nil(errListSingle)
		})
		t.Run("Create", func(t *testing.T) {
			m := manager.New("")
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
			m := manager.New("")
			_, errRmEmpty := m.Remove("")
			check.Nil(errRmEmpty)
			m.Create("test", "", false)
			_, errRmAll := m.Remove("")
			check.Nil(errRmAll)
			check.Equal(len(m.Data.Records), 0)
		})
		t.Run("Reset", func(t *testing.T) {
			m := manager.New("")
			_, err := m.Reset()
			check.Nil(err)
		})
	})
}
