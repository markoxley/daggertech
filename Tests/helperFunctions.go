package daggertests

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/markoxley/daggertech"
	"github.com/markoxley/daggertech/utils"
)

// TestModel for testing database
type TestModel struct {
	daggertech.Model

	Age   int    `daggertech:""`
	Name  string `daggertech:"size:20"`
	Death *int   `daggertech:""`
}

func getConnectionDetails() *daggertech.Config {
	c := daggertech.CreateConfig("tcp(127.0.0.1:3306)", "daggertechtest", "tester", "tester", true)
	return c
}

func getConnection() (*sql.DB, bool) {
	cs := "root:gbjbamox@tcp(127.0.0.1:3306)/daggertechtest?charset-utf8"
	if tdb, err := sql.Open("mysql", cs); err == nil {
		return tdb, true
	}
	return nil, false
}

func closeConnection(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func testTableExists(t string) bool {
	if c, ok := getConnection(); ok {
		defer closeConnection(c)
		sql := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE  table_schema = 'public' AND table_name = '%s');`, t)
		if r, err := c.Query(sql); err == nil {
			if r.Next() {
				return true
			}
		}
		return false
	}
	return false
}

func configuredaggertech() {
	daggertech.Configure(getConnectionDetails())
}

func reset() {
	configuredaggertech()
	sql := "Delete from TestModel;"
	if c, ok := getConnection(); ok {
		defer closeConnection(c)
		c.Exec(sql)
	}
}

func compareDates(s time.Time, d time.Time) bool {
	d1 := utils.TimeToSQL(&s)
	d2 := utils.TimeToSQL(&d)

	return d1 == d2
}
