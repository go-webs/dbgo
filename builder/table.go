package builder

import (
	"errors"
	"fmt"
	"gitub.com/go-webs/dbgo/util"
	"reflect"
	"strings"
)

type TableBuilder struct {
	table  any
	as     string
	prefix string
}

func NewTableBuilder(prefix string) *TableBuilder {
	ts := TableBuilder{
		prefix: prefix,
	}
	return &ts
}

// Table select from table as a
// param
//
//	arg: table string | struct bindings
//	as: table name alias to short name
func (ts *TableBuilder) Table(table any, as ...string) *TableBuilder {
	ts.table = table
	if len(as) > 0 {
		ts.as = as[0]
	}
	return ts
}

func (ts *TableBuilder) BuildTable() (sqlSegment string, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(ts.table))
	switch rfv.Kind() {
	case reflect.String:
		sqlSegment = fmt.Sprintf("`%s%s`", ts.prefix, ts.table)
	case reflect.Struct:
		sqlSegment = ts.buildTableName(rfv.Type())
	case reflect.Slice:
		if rfv.Type().Elem().Kind() == reflect.Struct {
			sqlSegment = ts.buildTableName(rfv.Type().Elem())
		} else {
			err = errors.New("table param must be string or struct(slice) bind with 1 or 2 params")
		}
	default:
		err = errors.New("table must string | struct | slice")
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", sqlSegment, ts.as)), err
}

func (ts *TableBuilder) buildTableName(rft reflect.Type) (tab string) {
	if field, ok := rft.FieldByName("TableName"); ok {
		if field.Tag.Get("db") != "" {
			tab = field.Tag.Get("db")
		}
	}
	if tab == "" {
		tab = fmt.Sprintf("`%s%s`", ts.prefix, strings.ToLower(rft.Name()))
	} else {
		tab = util.BackQuotes(tab)
	}
	return
}
