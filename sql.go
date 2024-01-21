package dbgo

import "fmt"

func (db Database) BuildQuery() (sqls string, values []any, err error) {
	return db.ToSqlOnly(), values, err
}
func (db Database) ToSqlOnly() string {
	var values []any

	var distinct = db.distinct
	var fields, bindValuesSelect = db.BuildSelect()
	//var table = db.BuildTable()
	var tables, _ = db.TableBuilder.BuildTable()
	//var join = db.BuildJoin()
	var joins, bindValuesJoin, _ = db.BuildJoin()
	wheres, bindValuesWhere, _ := db.BuildWhere()
	//db.BuildW
	var group = ""
	var having = ""
	var orderBy = ""
	var limit = ""
	var offset = ""

	values = append(values, bindValuesSelect...)
	values = append(values, bindValuesJoin...)
	values = append(values, bindValuesWhere...)

	if wheres != "" {
		wheres = fmt.Sprintf("WHERE %s", wheres)
	}

	return NamedSprintf("SELECT :distinct :fields FROM :tables :joins :wheres :group :having :orderBy :limit :offset",
		distinct, fields, tables, joins, wheres, group, having, orderBy, limit, offset)
}
