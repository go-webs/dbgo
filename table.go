package dbgo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type tableStruct struct {
	table any
	as    string
}

// Table string | struct bindings
func (db Database) Table(table any, as ...string) Database {
	db.tableStruct = tableStruct{
		table: table,
	}
	if len(as) > 0 {
		db.tableStruct.as = as[0]
	}
	return db
}

// TableAs select from table as a
// param
//
//	arg: table string | table struct bindings
//	as: table name alias
func (db Database) TableAs(table any, as string) Database {
	db.tableStruct = tableStruct{
		table: table,
		as:    as,
	}
	return db
}
func (db Database) BuildTable() (tab string) {
	return strings.TrimSpace(fmt.Sprintf("%s %s", db.buildTable(db.tableStruct.table), db.tableStruct.as))
}

// buildTable name
func (db Database) buildTable(tables any) (tab string) {
	rfv := reflect.Indirect(reflect.ValueOf(tables))
	switch rfv.Kind() {
	case reflect.String:
		tab = fmt.Sprintf("`%s%s`", db.DbGo.Prefix, tables)
	case reflect.Struct:
		tab = db.buildTableName(rfv.Type())
	case reflect.Slice:
		if rfv.Type().Elem().Kind() == reflect.Struct {
			tab = db.buildTableName(rfv.Type().Elem())
		} else {
			db.Err = errors.New("table param must be string or struct(slice) bind with 1 or 2 params")
		}
		//if rfv.Len() == 0 { // *[]Users
		//	if rfv.Type().Elem().Kind() == reflect.Struct {
		//		tab = db.buildTable(reflect.New(rfv.Type().Elem()).Interface())
		//	}
		//} else if rfv.Len() == 1 {
		//	tab = db.buildTable(rfv.Index(0).Interface())
		//} else if rfv.Len() == 2 {
		//	tab = db.buildTable(
		//		rfv.Index(0).Interface())
		//	tab = fmt.Sprintf("%s %s", tab, rfv.Index(1))
		//} else {
		//	db.Err = errors.New("table param must be string or struct(slice) bind with 1 or 2 params")
		//}
	default:
		db.Err = errors.New("table must string | struct | slice")
	}
	if db.Err != nil {
		return
	}

	return
}

func (db Database) buildTableName(rft reflect.Type) (tab string) {
	if field, ok := rft.FieldByName("TableName"); ok {
		if field.Tag.Get("db") != "" {
			tab = field.Tag.Get("db")
		}
	}
	if tab == "" {
		tab = fmt.Sprintf("`%s%s`", db.DbGo.Prefix, rft.Name())
	} else {
		tab = fmt.Sprintf("`%s`", tab)
	}
	return
}
