package dbgo

// OrderBy clause
func (db Database) OrderBy(args ...string) Database {
	db.OrderByBuilder.OrderBy(args...)
	return db
}

// OrderByAsc clause
func (db Database) OrderByAsc(arg string) Database {
	db.OrderBy(arg, "ASC")
	return db
}

// OrderByDesc clause
func (db Database) OrderByDesc(arg string) Database {
	db.OrderBy(arg, "DESC")
	return db
}

// OrderByRaw clause
func (db Database) OrderByRaw(args ...string) Database {
	db.OrderByBuilder.OrderByRaw(args...)
	return db
}

// Page clause
func (db Database) Page(arg int) Database {
	db.PageBuilder.Page(arg)
	return db
}

// Limit clause
func (db Database) Limit(arg int) Database {
	db.PageBuilder.Limit(arg)
	return db
}

// Offset clause
func (db Database) Offset(arg int) Database {
	db.PageBuilder.Offset(arg)
	return db
}

// Take clause
func (db Database) Take(arg int) Database {
	db.PageBuilder.Limit(arg)
	return db
}

// Skip clause
func (db Database) Skip(arg int) Database {
	db.PageBuilder.Offset(arg)
	return db
}
