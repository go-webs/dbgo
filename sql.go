package dbgo

import (
	"fmt"
	"gitub.com/go-webs/dbgo/util"
	"reflect"
	"strings"
)

func (db Database) BuildSqlQuery() (sqlSegment string, values []any, err error) {
	var distinct = db.distinct
	var fields, bindValuesSelect = db.BuildSelect()
	tables, err := db.TableBuilder.BuildTable()
	if err != nil {
		return sqlSegment, values, err
	}
	joins, bindValuesJoin, err := db.BuildJoin()
	if err != nil {
		return sqlSegment, values, err
	}
	wheres, bindValuesWhere, err := db.BuildWhere()
	if err != nil {
		return sqlSegment, values, err
	}
	groups, havingS, bindValuesGroup := db.BuildGroup()
	if err != nil {
		return sqlSegment, values, err
	}
	orderBys := db.BuildOrderBy()
	pagination, bindValuesPagination := db.BuildPage()

	values = append(values, bindValuesSelect...)
	values = append(values, bindValuesJoin...)
	values = append(values, bindValuesWhere...)
	values = append(values, bindValuesGroup...)
	values = append(values, bindValuesPagination...)

	if wheres != "" {
		wheres = fmt.Sprintf("WHERE %s", wheres)
	}
	//else {
	//	// 从struct构建where
	//}

	sqlSegment = util.NamedSprintf("SELECT :distinct :fields FROM :tables :joins :wheres :groups :havings :orderBys :page",
		distinct, fields, tables, joins, wheres, groups, havingS, orderBys, pagination)
	return
}
func (db Database) BuildSqlInsert(data any) (sqlSegment string, values []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(data))
	//var fn = func(dataRfv reflect.Value) {
	//	var valueTmp []any
	//	valueTmp = append(valueTmp, rfv.MapIndex(key).Interface())
	//	keys := rfv.MapKeys()
	//	for _, key := range keys {
	//		valueTmp = append(valueTmp, rfv.MapIndex(key).Interface())
	//	}
	//}
	var fields []string
	var valuesPlaceholderArr []string
	switch rfv.Kind() {
	case reflect.Map:
		keys := rfv.MapKeys()
		var valuesPlaceholderTmp []string
		for _, key := range keys {
			fields = append(fields, util.BackQuotes(key.String()))
			valuesPlaceholderTmp = append(valuesPlaceholderTmp, "?")
			values = append(values, rfv.MapIndex(key).Interface())
		}
		valuesPlaceholderArr = append(valuesPlaceholderArr, fmt.Sprintf("(%s)", strings.Join(valuesPlaceholderTmp, ",")))
	case reflect.Slice:
		if rfv.Len() == 0 {
			return
		}
		// 先获取到插入字段
		keys := rfv.Index(0).MapKeys()
		for _, key := range keys {
			fields = append(fields, util.BackQuotes(key.String()))
		}
		// 组合插入数据
		for i := 0; i < rfv.Len(); i++ {
			var valuesPlaceholderTmp []string
			for _, key := range keys {
				valuesPlaceholderTmp = append(valuesPlaceholderTmp, "?")
				values = append(values, rfv.Index(i).MapIndex(key).Interface())
			}
			valuesPlaceholderArr = append(valuesPlaceholderArr, fmt.Sprintf("(%s)", strings.Join(valuesPlaceholderTmp, ",")))
		}
	}
	tables := db.BuildTable()
	sqlSegment = util.NamedSprintf("INSERT INTO :tables (:fields) VALUES :placeholder", tables, strings.Join(fields, ","), strings.Join(valuesPlaceholderArr, " "))
	return
}
func (db Database) BuildSqlUpdate() (sqlSegment string, values []any, err error) { return }
func (db Database) BuildSqlDelete() (sqlSegment string, values []any, err error) { return }
func (db Database) ToSqlOnly() string {
	sqls, _, _ := db.BuildSqlQuery()
	return sqls
}
