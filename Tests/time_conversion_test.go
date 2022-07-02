package daggertests

import (
	"testing"
	"time"

	"github.com/markoxley/daggertech/utils"
)

func TestValidTimeToUint64AndBack(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 24, 0, time.UTC)
	e := utils.TimeToUint64(&tm1)
	tm2, ok := utils.Uint64ToTime(e)
	if !ok {
		t.Error("Unable to convert time to Uint64")
	}
	if !compareDates(tm1, *tm2) {
		t.Errorf("Expected %v, got %v", tm1, tm2)
	}
}

func TestFromUint64ToTimeAndBack(t *testing.T) {
	e := uint64(63191276532652)
	tm, ok := utils.Uint64ToTime(e)
	if !ok {
		t.Error("First conversion not OK")
	}
	i := utils.TimeToUint64(tm)
	if i != e {
		t.Errorf("Starting value is %v, end value is %v", e, i)
	}
}

func TestFromSQLToTimeAndBack(t *testing.T) {
	org := "2020-09-18 10:35:34.786"
	tm, ok := utils.SQLToTime(org)
	if !ok {
		t.Error("Invalid SQL time format")
	}
	res := utils.TimeToSQL(tm)
	if res != org {
		t.Errorf("Starting value was '%s', end value is '%s'", org, res)
	}
}

func TestFromTimeToSQLAndBack(t *testing.T) {
	org := time.Now()
	sql := utils.TimeToSQL(&org)
	res, ok := utils.SQLToTime(sql)
	if !ok {
		t.Error("Invalid SQL time format")
	}
	if !compareDates(org, *res) {
		t.Errorf("Expecting '%s', got '%s'", utils.TimeToSQL(&org), utils.TimeToSQL(res))
	}
}
