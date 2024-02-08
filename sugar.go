package dbgo

func (db Database) MustFirst(columns ...string) (res map[string]any) {
	res, db.Context.Err = db.First(columns...)
	return
}
func (db Database) MustGet(columns ...string) (res []map[string]any) {
	res, db.Context.Err = db.Get(columns...)
	return
}
func (db Database) WhereIn(column string, value any) Database {
	db.Context.WhereClause.WhereIn("AND", column, value)
	return db
}
func (db Database) WhereNotIn(column string, value any) Database {
	db.Context.WhereClause.WhereIn("AND", column, value, true)
	return db
}
func (db Database) OrWhereIn(column string, value any) Database {
	db.Context.WhereClause.WhereIn("OR", column, value)
	return db
}
func (db Database) OrWhereNotIn(column string, value any) Database {
	db.Context.WhereClause.WhereIn("OR", column, value, true)
	return db
}
func (db Database) WhereNull(column string) Database {
	db.Context.WhereClause.WhereNull("AND", column)
	return db
}
func (db Database) WhereNotNull(column string) Database {
	db.Context.WhereClause.WhereNull("AND", column, true)
	return db
}
func (db Database) OrWhereNull(column string) Database {
	db.Context.WhereClause.WhereNull("OR", column)
	return db
}
func (db Database) OrWhereNotNull(column string) Database {
	db.Context.WhereClause.WhereNull("OR", column, true)
	return db
}
func (db Database) WhereBetween(column string, value any) Database {
	db.Context.WhereClause.WhereBetween("AND", column, value)
	return db
}
func (db Database) WhereNotBetween(column string, value any) Database {
	db.Context.WhereClause.WhereBetween("AND", column, value, true)
	return db
}
func (db Database) OrWhereBetween(column string, value any) Database {
	db.Context.WhereClause.WhereBetween("OR", column, value)
	return db
}
func (db Database) OrWhereNotBetween(column string, value any) Database {
	db.Context.WhereClause.WhereBetween("OR", column, value, true)
	return db
}
func (db Database) WhereLike(column, value string) Database {
	db.Context.WhereClause.WhereLike("AND", column, value)
	return db
}
func (db Database) WhereNotLike(column, value string) Database {
	db.Context.WhereClause.WhereLike("AND", column, value, true)
	return db
}
func (db Database) OrWhereLike(column, value string) Database {
	db.Context.WhereClause.WhereLike("OR", column, value)
	return db
}
func (db Database) OrWhereNotLike(column, value string) Database {
	db.Context.WhereClause.WhereLike("OR", column, value, true)
	return db
}

func (db Database) OrderByAsc(column string) Database {
	return db.OrderBy(column, "ASC")
}

func (db Database) OrderByDesc(column string) Database {
	return db.OrderBy(column, "DESC")
}
