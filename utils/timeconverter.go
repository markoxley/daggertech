package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeToUint64 returns the date and time in numerical format for storing in a database
func TimeToUint64(t *time.Time) uint64 {

	y := uint64(t.Year())
	m := uint64(t.Month())
	d := uint64(t.Day())
	h := uint64(t.Hour())
	mn := uint64(t.Minute())
	s := uint64(t.Second())
	n := uint64(t.Nanosecond())

	r := (((((((((((y * 12) + m) * 31) + d) * 24) + h) * 60) + mn) * 60) + s) * 1000) + (n / 1000000)

	return r
}

// Uint64ToTime converts the numeric version of the date and time stored in the database back to a Time struct
func Uint64ToTime(i uint64) (*time.Time, bool) {
	ok := true

	defer func() {
		if x := recover(); x != nil {
			ok = false
		}
	}()

	var y, m, d, h, mn, s, ns int
	ns = int(i%1000) * 1000000
	i /= 1000
	s = int(i % 60)
	i /= 60
	mn = int(i % 60)
	i /= 60
	h = int(i % 24)
	i /= 24
	d = int(i % 31)
	i /= 31
	m = int(i % 12)
	y = int(i / 12)

	t := time.Date(y, time.Month(m), d, h, mn, s, ns, time.UTC)
	return &t, ok
}

// SQLToTime attempts to convert a string, containing an SQL datetime string to a time object
func SQLToTime(st string) (*time.Time, bool) {
	sep := " "
	var y, m, d, h, mn, s, ns int
	var e error
	if strings.Contains(st, "T") {
		sep = "T"
	}
	dt := strings.Split(st, sep)
	if len(dt) < 1 {
		return nil, false
	}
	dp := strings.Split(dt[0], "-")
	if len(dp) != 3 {
		return nil, false
	}
	if y, e = strconv.Atoi(dp[0]); e != nil {
		return nil, false
	}
	if m, e = strconv.Atoi(dp[1]); e != nil {
		return nil, false
	}
	if d, e = strconv.Atoi(dp[2]); e != nil {
		return nil, false
	}
	if len(dt) > 1 {
		tm := strings.Split(dt[1], ":")
		if len(tm) != 3 {
			return nil, false
		}

		sc := strings.Split(tm[2], ".")

		if len(sc) > 1 {
			if ns, e = strconv.Atoi(sc[1]); e != nil {
				ns = 0
			}
		}
		if h, e = strconv.Atoi(tm[0]); e != nil {
			return nil, false
		}
		if mn, e = strconv.Atoi(tm[1]); e != nil {
			return nil, false
		}
		if s, e = strconv.Atoi(sc[0]); e != nil {
			return nil, false
		}
	}
	t := time.Date(y, time.Month(m), d, h, mn, s, ns, time.UTC)
	return &t, true
}

// TimeToSQL formats the time value to an SQL string
func TimeToSQL(t *time.Time) string {
	var y, m, d, h, mn, s, ns int
	y = t.Year()
	m = int(t.Month())
	d = t.Day()
	h = t.Hour()
	mn = t.Minute()
	s = t.Second()
	ns = t.Nanosecond()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%d", y, m, d, h, mn, s, ns)
}
