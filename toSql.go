package dbgo2

func (db Database) ToSql() (sql4prepare string, values []any, err error) {
	return db.Driver().ToSql(db.Context)
}

func (db Database) ToSqlSelect() (sql4prepare string, binds []any) {
	return db.Driver().ToSqlSelect(db.Context)
}

func (db Database) ToSqlTable() (sql4prepare string, values []any, err error) {
	return db.Driver().ToSqlTable(db.Context)
}

func (db Database) ToSqlWhere() (sql4prepare string, values []any, err error) {
	return db.Driver().ToSqlWhere(db.Context)
}

func (db Database) ToSqlOrderBy() (sql4prepare string) {
	return db.Driver().ToSqlOrderBy(db.Context)
}

func (db Database) ToSqlLimitOffset() (sqlSegment string, binds []any) {
	return db.Driver().ToSqlLimitOffset(db.Context)
}
