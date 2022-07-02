package daggertech

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/markoxley/daggertech/utils"
)

type fieldMap struct {
	name      string
	fieldType string
}

var (
	conf        *Config
	knownTables []string
	//gdb  *sql.DB
)

const (
	// dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"
	// db, err := sql.Open("mysql", dsn)
	connectionPattern = "%s:%s@%s/%s"
	// connectionPattern = "host=%s port=%d dbname=%s user=%s password=%s sslmode=disable"
)

func init() {
	knownTables = make([]string, 0, 20)
}

// Configure loads the configraution for connection to the database.
// This must be done before all other database operations
func Configure(c *Config) bool {
	conf = c
	db, err := connect()
	if err != nil {
		return false
	}
	db.Close()
	return true
}

// Connect to the database
func connect() (*sql.DB, error) {
	// if gdb != nil {
	// 	return gdb, nil
	// }
	//cs := fmt.Sprintf(connectionPattern, conf.host, conf.port, conf.name, conf.user, conf.password)
	cs := fmt.Sprintf(connectionPattern, conf.user, conf.password, conf.host, conf.name)
	tdb, err := sql.Open("mysql", cs)
	if err != nil {
		return nil, err
	}
	//gdb = tdb
	return tdb, nil
}

// Disconnect from the database
func disconnect(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func beginTransaction(db *sql.DB) (*sql.Tx, error) {
	return db.Begin()
}

func commitTransaction(tx *sql.Tx) {
	if tx != nil {
		tx.Commit()
	}
}

func selectScalar(q string) (interface{}, bool) {
	db, err := connect()
	if err != nil {
		return nil, false
	}
	defer disconnect(db)

	tx, err := beginTransaction(db)
	if err != nil {
		return nil, false
	}
	defer commitTransaction(tx)

	res, err := db.Query(q)
	if err != nil {
		return nil, false
	}
	defer res.Close()
	if res.Next() {
		var cols string
		vl := &cols
		//var vl interface{}
		res.Scan(vl)
		return cols, true
	}
	return nil, false

}

// Perform a sellect query and return all the rows
func selectQuery(q string, m Modeller) ([]Modeller, bool) {
	db, err := connect()
	if err != nil {
		return nil, false
	}
	defer disconnect(db)

	tx, err := beginTransaction(db)
	if err != nil {
		return nil, false
	}
	defer commitTransaction(tx)

	res, err := db.Query(q)
	if err != nil {
		return nil, false
	}
	defer res.Close()
	return populateModel(m, res)
}

func populateModel(m Modeller, r *sql.Rows) ([]Modeller, bool) {
	res := make([]Modeller, 0, 10)
	ok := true
	// Get the column count
	cc, _ := r.Columns()

	// Make them all uppercase
	for i := range cc {
		cc[i] = strings.ToUpper(cc[i])
	}

	// Get the fields of the model and build a map of them
	//t := reflect.TypeOf(*m)
	flds, ok := tableDef[getTableName(m)]
	if !ok {
		return nil, false
	}
	fMap := make(map[string]field, len(flds))
	for _, f := range flds {
		fMap[strings.ToUpper(f.name)] = f
	}

	cols := make([]*string, len(cc))
	vls := make([]interface{}, len(cc))

	for r.Next() {

		s := reflect.ValueOf(m).Elem().Type()
		v := reflect.New(s)

		for i := range cols {
			vls[i] = &cols[i]
		}
		r.Scan(vls...)
		tmpID := ""
		tmpCreate := time.Now()
		tmpUpdate := time.Now()
		var tmpDeleted *time.Time

		for i := 0; i < len(cc); i++ {
			if cols[i] == nil {
				continue
			}
			if cc[i] == "ID" {
				tmpID = *cols[i]
			} else if cc[i] == "CREATEDATE" {
				if cols[i] != nil {
					if tm, ok := utils.SQLToTime(*cols[i]); ok {
						tmpCreate = *tm
					}
				}
				// if val, err := strconv.ParseUint(*cols[i], 10, 64); err == nil {
				// 	if tm, ok := utils.Uint64ToTime(val); ok {
				// 		tmpCreate = *tm
				// 	}
				// }
			} else if cc[i] == "LASTUPDATE" {
				if cols[i] != nil {
					if tm, ok := utils.SQLToTime(*cols[i]); ok {
						tmpUpdate = *tm
					}
					// if val, err := strconv.ParseUint(*cols[i], 10, 64); err == nil {
					// 	if tm, ok := utils.Uint64ToTime(val); ok {
					// tmpUpdate = *tm
					// }
					// }
				}
			} else if cc[i] == "DELETEDATE" {
				if cols[i] != nil {
					if tm, ok := utils.SQLToTime(*cols[i]); ok {
						tmpDeleted = tm
					}
					// if val, err := strconv.ParseUint(*cols[i], 10, 64); err == nil {
					// 	if tm, ok := utils.Uint64ToTime(val); ok {
					// 		tmpDeleted = tm
					// 	}
					// }
				}
			} else if fld, ok := fMap[cc[i]]; ok {
				switch fld.fType {
				case tInt, tLong:
					if fld.unsigned {
						if val, err := strconv.ParseUint(*cols[i], 10, 64); err != nil {
							if fld.allowNull {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
							} else {
								v.Elem().FieldByName(fld.name).SetUint(val)
							}
						}
					} else {
						if val, err := strconv.ParseInt(*cols[i], 10, 64); err == nil {
							if fld.allowNull {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
							} else {
								v.Elem().FieldByName(fld.name).SetInt(val)
							}
						}
					}
				case tBool:
					if val, err := strconv.ParseInt(*cols[i], 10, 0); err == nil {
						if fld.allowNull {
							//boolVal := val == 1
							v.Elem().FieldByName(fld.name).Elem().SetBool(val == 1) //Set(reflect.ValueOf(&boolVal))
						} else {
							v.Elem().FieldByName(fld.name).SetBool(val == 1)
						}
					}
				case tDecimal, tFloat, tDouble:
					if val, err := strconv.ParseFloat(*cols[i], 64); err == nil {
						if fld.allowNull {
							v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
						} else {
							v.Elem().FieldByName(fld.name).SetFloat(val)
						}
					}
				case tDateTime:
					if cols[i] != nil {
						if val, ok := utils.SQLToTime(*cols[i]); ok {
							if fld.allowNull {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
							} else {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(*val))
							}
						}
					}
					// if val, err := strconv.ParseUint(*cols[i], 10, 64); err == nil {
					// 	if tm, ok := utils.Uint64ToTime(val); ok {
					// 		if fld.allowNull {
					// 			v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(tm))
					// 		} else {
					// 			v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(*tm))
					// 		}
					// 	}
					// }
				case tChar:
					st := *cols[i]
					strVal := st[:1]
					if fld.allowNull {
						v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
					} else {
						v.Elem().FieldByName(fld.name).SetString(strVal)
					}
				case tString:
					if fld.allowNull {
						strVal := *cols[i]
						v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
					} else {
						v.Elem().FieldByName(fld.name).SetString(*cols[i])
					}
				case tUUID:
					if fld.allowNull {
						strVal := *cols[i]
						v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
					} else {
						v.Elem().FieldByName(fld.name).SetString(*cols[i])
					}
				}
			}

		}
		newObj := v.Interface().(Modeller)
		updateModel(newObj, tmpID, tmpCreate, tmpUpdate, tmpDeleted)
		res = append(res, newObj)

	}

	return res, ok
}

func executeQuery(q string) bool {
	db, err := connect()
	if err != nil {
		return false
	}
	defer disconnect(db)

	tx, err := beginTransaction(db)
	if err != nil {
		return false
	}
	defer commitTransaction(tx)

	_, err = db.Exec(q)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func tableExists(t string) bool {
	for _, tn := range knownTables {
		if tn == t {
			return true
		}
	}
	if _, ok := selectScalar(fmt.Sprintf("SHOW TABLES WHERE Tables_in_%s = '%s'", conf.name, t)); ok {
		knownTables = append(knownTables, t)
		return true
	}
	return false
}

// RawExecute executes a sql statement on the database, without returning a value
// Not recommended for general use - can break shadowing
func RawExecute(sql string) bool {
	return executeQuery(sql)
}

// RawScalar exeutes a raw sql statement that returns a single value
// Not recommended for general use
func RawScalar(sql string) (interface{}, bool) {
	return selectScalar(sql)
}

// RawSelect executes a raw sql statement on the database
// Not recommended for general use
func RawSelect(sql string) []map[string]interface{} {
	db, err := connect()
	if err != nil {
		return nil
	}
	defer disconnect(db)
	res, err := db.Query(sql)
	if err != nil {
		return nil
	}
	defer res.Close()
	data := make([]map[string]interface{}, 0, 10)

	// Get the column count
	cc, _ := res.Columns()

	cols := make([]*string, len(cc))
	vls := make([]interface{}, len(cc))

	for res.Next() {

		for i := range cols {
			vls[i] = &cols[i]
		}
		res.Scan(vls...)
		row := make(map[string]interface{})
		for i, n := range cc {
			row[n] = vls[i]
		}
		data = append(data, row)
	}
	return data
}

// Fetch populates the slice with models from the database that match the criteria.
// Returns false if this fails
func Fetch(m Modeller, c *Criteria) ([]Modeller, bool) {
	t := reflect.TypeOf(m)
	n := t.Name()
	//	nmNew := reflect.New(t).Elem().Interface()
	//	nm, _ := nmNew.(Modeller)
	_, n, ok := tableTest(m)
	if !ok {
		return nil, false
	}
	s := fmt.Sprintf("select * from `%s`", n)
	whereDone := false
	if c != nil {
		if c.Where != "" {
			s += fmt.Sprintf(" WHERE %s", c.Where)
			whereDone = true
		}
		if !c.IncDeleted {
			if whereDone {
				s += " AND"
			} else {
				s += "WHERE"
			}
			s += "`DeleteDate` Is Null"
		}

		if c.Order != "" {
			s += fmt.Sprintf(" ORDER BY %s", c.Order)
		}
		if c.Limit > 0 {
			s += fmt.Sprintf(" LIMIT %d", c.Limit)
		}
		if c.Offset > 0 {
			s += fmt.Sprintf(" OFFSET %d", c.Offset)
		}
	}
	res, ok := selectQuery(s, m)
	return res, ok
}

// First populates the model with the first record from the database that match the criteria.
// Returns false if this fails or there is no record to return
func First(m Modeller, c *Criteria) (Modeller, bool) {
	if c == nil {
		c = &Criteria{}
	}
	c.Limit = 1
	c.Offset = 0
	if ml, ok := Fetch(m, c); ok {
		if len(ml) > 0 {
			return ml[0], true
		}
	}
	return nil, false
}

// Count returns the number of records in the named table that match the criteria
func Count(t string, c *Criteria) int {
	if !tableExists(t) {
		return 0
	}
	s := fmt.Sprintf("Select Count(*) from `%s`", t)
	whereAdded := false
	if c != nil {
		if c.Where != "" {
			s += fmt.Sprintf(" WHERE %s", c.Where)
			whereAdded = true
		}
	}
	if c == nil || !c.IncDeleted {
		if whereAdded {
			s += " AND DeleteDate is null"
		} else {
			s += " WHERE DeleteDate is null"
		}
	}
	if i, ok := selectScalar(s); ok {
		if vl, vlok := i.(string); vlok {
			if res, err := strconv.Atoi(vl); err == nil {
				return res
			}
		}
	}
	return 0

}

// Save stores the model in the database
func Save(m Modeller) bool {
	if m.IsNew() {
		return executeQuery(insertCommand(m))
	}
	return executeQuery(updateCommand(m))
}

// Remove deletes the model from the database
func Remove(m Modeller) bool {
	if m.GetID() == nil {
		return true
	}
	if conf.deletable {
		return executeQuery(fmt.Sprintf("delete from `%s` where id = '%s'", getTableName(m), *(m.GetID())))
	}
	now := time.Now()
	return executeQuery(fmt.Sprintf("update `%s` set `deleteDate` = %v where `id` = '%s'", getTableName(m), utils.TimeToSQL(&now), *(m.GetID())))
}

// RemoveMany removes all records from the named table that match the criteria
func RemoveMany(t string, c *Criteria) (int, bool) {
	if !tableExists(t) {
		return 0, true
	}
	r := Count(t, c)
	if r == 0 {
		return 0, true
	}
	s := ""
	if conf.deletable {
		s = fmt.Sprintf("delete from %s", t)
	} else {
		tm := time.Now()
		s = fmt.Sprintf("update %s set `deleteDate` = '%v'", t, utils.TimeToSQL(&tm))
	}
	whereAdded := false
	if c != nil && c.Where != "" {
		s += fmt.Sprintf(" where %s", c.Where)
		whereAdded = true
	}

	if whereAdded {
		s += " AND DeleteDate is null"
	} else {
		s += " WHERE DeleteDate is null"
	}
	// if !conf.deletable {
	// 	createShadow()
	// }
	ok := executeQuery(s)
	return r, ok
}
