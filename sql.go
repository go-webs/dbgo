package dbgo

import (
	"fmt"
	"gitub.com/go-webs/dbgo/util"
)

func (db Database) BuildSqlQuery() (sqlSegment string, values []any, err error) {
	return db.ToSqlOnly(), values, err
}
func (db Database) BuildSqlInsert() (sqlSegment string, values []any, err error) { return }
func (db Database) BuildSqlUpdate() (sqlSegment string, values []any, err error) { return }
func (db Database) BuildSqlDelete() (sqlSegment string, values []any, err error) { return }
func (db Database) ToSqlOnly() string {
	var values []any

	var distinct = db.distinct
	var fields, bindValuesSelect = db.BuildSelect()
	var tables, _ = db.TableBuilder.BuildTable()
	joins, bindValuesJoin, _ := db.BuildJoin()
	wheres, bindValuesWhere, _ := db.BuildWhere()
	groups, havingS, _ := db.BuildGroup()
	orderBys := db.BuildOrderBy()
	pagination, bindValuesPagination := db.BuildPage()

	values = append(values, bindValuesSelect...)
	values = append(values, bindValuesJoin...)
	values = append(values, bindValuesWhere...)
	values = append(values, bindValuesPagination...)

	if wheres != "" {
		wheres = fmt.Sprintf("WHERE %s", wheres)
	} else {
		// 从struct构建where
	}

	return util.NamedSprintf("SELECT :distinct :fields FROM :tables :joins :wheres :groups :havings :orderBys :page",
		distinct, fields, tables, joins, wheres, groups, havingS, orderBys, pagination)
}
