package daggertech

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/markoxley/daggertech/clause"
	"github.com/markoxley/daggertech/utils"
	uuid "github.com/satori/go.uuid"
)

// ModelState is the state of the current model
type ModelState int

// Model is the base for all database models
type Model struct {
	ID         *string
	CreateDate time.Time
	LastUpdate time.Time
	DeleteDate *time.Time
	tableName  *string
}

// CreateModel sets the default parameters for the Model
func CreateModel() Model {
	return Model{
		CreateDate: time.Now(),
		LastUpdate: time.Now(),
	}
}

// StandingData returns the standing data for when the table is created
func (m Model) StandingData() []Modeller {
	return nil
}

// GetID returns the ID of the model
func (m Model) GetID() *string {
	return m.ID
}

// IsNew returns true if the model has yet to be stored
func (m Model) IsNew() bool {
	return m.ID == nil
}

// IsDeleted returns true if teh model has been marked as deleted
func (m Model) IsDeleted() bool {
	return m.DeleteDate == nil
}

func getTableName(m Modeller) string {
	return reflect.Indirect(reflect.ValueOf(m).Elem()).Type().Name()
}

func tableTest(m Modeller) ([]field, string, bool) {
	sql, required := tableDefinition(m)
	if required {
		te := tableExists(getTableName(m))
		knownTables = append(knownTables, getTableName(m))
		if !te {
			for _, s := range sql {
				if !executeQuery(s) {
					log.Panicf(`Error executing "%s"`, s)
				}
			}
			if standingData := m.StandingData(); standingData != nil {
				for _, data := range standingData {
					Save(data)
				}
			}
		}
	}
	flds, ok := tableDef[getTableName(m)]
	return flds, getTableName(m), ok
}

// Returns a slice of strings with the sql statements and boolean to indicate if the table needs to be created
func tableDefinition(m Modeller) ([]string, bool) {
	sql := make([]string, 0, 3)

	n := getTableName(m)
	if _, ok := tableDef[n]; ok {
		return nil, false
	}

	t := reflect.TypeOf(m).Elem()
	nm := reflect.New(t).Elem().Interface()

	fs := getDefs(nm, true)

	tableDef[n] = fs
	if len(fs) == 0 {
		return nil, false
	}
	//flds := "id varchar(36) primary key, createDate bigint, lastUpdate bigint, disabled tinyint"
	flds := ""
	keys := make([]string, 0, 5)
	for _, f := range fs {
		if flds != "" {
			flds += ", "
		}
		flds += fmt.Sprintf("`%s` %s", f.name, pgFieldNames[f.fType])
		if f.fType != tUUID && f.fType != tChar && f.size.size > 0 {
			flds += fmt.Sprintf("(%s)", f.size.String())
		}
		if f.fType == tString && f.size.size == 0 {
			flds += "(256)"
		}
		if f.unsigned {
			flds += " UNSIGNED"
		}
		if !f.allowNull {
			flds += " NOT NULL"
		}
		if f.key {
			keys = append(keys, f.name)
		}
	}
	sql = append(sql, fmt.Sprintf(sqlTableCreate, n, flds))
	kn := strings.ReplaceAll(n, ".", "_")
	// sql = append(sql, fmt.Sprintf(sqlIndexCreate, kn, "createDate", n, "createDate"))
	// sql = append(sql, fmt.Sprintf(sqlIndexCreate, kn, "lastUpdate", n, "lastUpdate"))
	for _, k := range keys {
		sql = append(sql, fmt.Sprintf(sqlIndexCreate, kn, k, n, k))
	}
	return sql, true
}

func insertCommand(m Modeller) string {
	flds, n, ok := tableTest(m)
	if !ok {
		return ""
	}
	uid := uuid.NewV4()

	fds := "ID, CreateDate, LastUpdate"
	now := time.Now()
	dbNow := utils.TimeToSQL(&now)
	updateModel(m, fmt.Sprintf("%s", uid), now, now, nil)
	q := fmt.Sprintf("'%s', '%s', '%s'", uid, dbNow, dbNow)
	v := reflect.ValueOf(m).Elem()
	//vt := reflect.TypeOf(m).Elem()
	for _, f := range flds {
		if f.name == "ID" || f.name == "CreateDate" || f.name == "LastUpdate" || f.name == "DeleteDate" {
			continue
		}
		vi := v.FieldByName(f.name)

		if f.allowNull {
			if vi.IsNil() {
				continue
			}
			vi = vi.Elem()
		}

		vf := vi.Interface()

		if vl, _, ok := clause.MakeValue(vf); ok {
			fds += fmt.Sprintf(", `%s`", f.name)
			q += fmt.Sprintf(", %s", vl)
		}
	}

	def := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", n, fds, q)
	return def
}

func updateCommand(m Modeller) string {
	flds, n, ok := tableTest(m)
	if !ok {
		return ""
	}
	now := time.Now()
	updateLastUpdate(m, now)
	res := fmt.Sprintf("UPDATE %s SET", n)
	v := reflect.ValueOf(m)
	first := true
	for _, f := range flds {
		if f.name != "ID" && f.name != "CreateDate" {
			if !first {
				res += ","
			}
			first = false
			var value interface{}
			if f.allowNull {
				if v.Elem().FieldByName(f.name).IsNil() {
					res += fmt.Sprintf(" `%s` = null", f.name)
					continue
				}
				value = v.Elem().FieldByName(f.name).Elem().Interface()
			} else {
				value = v.Elem().FieldByName(f.name).Interface()
			}
			if vl, _, ok := clause.MakeValue(value); ok {
				res += fmt.Sprintf(" `%s` = %s", f.name, vl)
			}
		}
	}
	def := res + fmt.Sprintf(" WHERE `Id` = '%s'", *m.GetID())
	return def
}

func deleteCommand(m Modeller) string {
	_, n, ok := tableTest(m)
	if !ok {
		return ""
	}
	def := fmt.Sprintf("DELETE FROM %s WHERE `Id` = '%s'", n, *m.GetID())
	return def
}

func refreshCommand(m Modeller) string {
	_, n, ok := tableTest(m)
	if !ok {
		return ""
	}
	def := fmt.Sprintf("SELECT * FROM %s WHERE `Id` = '%s'", n, *m.GetID())
	return def
}

func updateModel(m Modeller, id string, createdate time.Time, updatedate time.Time, deletedate *time.Time) {
	v := reflect.ValueOf(m)
	createdateValue := reflect.ValueOf(createdate)
	updatedateValue := reflect.ValueOf(updatedate)
	deletedateValue := reflect.ValueOf(deletedate)
	rv := reflect.New(reflect.TypeOf(id))
	rv.Elem().Set(reflect.ValueOf(id))

	v.Elem().FieldByName("ID").Set(rv)
	v.Elem().FieldByName("CreateDate").Set(createdateValue)
	v.Elem().FieldByName("LastUpdate").Set(updatedateValue)
	v.Elem().FieldByName("DeleteDate").Set(deletedateValue)
}

func updateLastUpdate(m Modeller, date time.Time) {
	v := reflect.ValueOf(m)
	dateValue := reflect.ValueOf(date)
	v.Elem().FieldByName("LastUpdate").Set(dateValue)
}
