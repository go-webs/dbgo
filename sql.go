package dbgo

func (db Database) ToSql() string {
	var values []any

	var distinct = db.distinct
	var fields, binds = db.BuildSelect()
	var table = db.buildTable(db.tables)
	var join = db.BuildJoin()
	var where = ""
	var group = ""
	var having = ""
	var orderBy = ""
	var limit = ""
	var offset = ""

	values = append(values, binds...)

	return NamedSprintf("SELECT :distinct :fields FROM :table :join :where :group :having :orderBy :limit :offset",
		distinct, fields, table, join, where, group, having, orderBy, limit, offset)
}
