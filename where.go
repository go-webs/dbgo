package dbgo

import "fmt"

// Where : query or execute where condition, the relation is and
func (db Database) Where(args ...interface{}) Database {
	if len(args) == 0 {
		return db
	}
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组
	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	//w := []interface{}{"AND", args}
	//db.where = append(db.where, w)
	db.WhereBuilder.Where(args...)
	return db
}

// OrWhere : query or execute where condition, the relation is or
func (db Database) OrWhere(args ...interface{}) Database {
	if len(args) == 0 {
		return db
	}
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组
	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	//w := []interface{}{"OR", args}
	//db.where = append(db.where, w)
	db.WhereBuilder.OrWhere(args...)

	return db
}

// WhereRaw fields with binds
// params
//
//	arg: expressions
//	binds: bind values
func (db Database) WhereRaw(arg string, binds ...any) Database {
	db.WhereBuilder.WhereRaw(arg, binds)
	return db
}

// OrWhereRaw fields with binds
func (db Database) OrWhereRaw(arg string, binds ...any) Database {
	db.WhereBuilder.OrWhereRaw(arg, binds)
	return db
}

// WhereNull ...
func (db Database) WhereNull(arg string) Database {
	return db.Where(fmt.Sprintf("`%s` IS NULL", arg))
}

// OrWhereNull ...
func (db Database) OrWhereNull(arg string) Database {
	return db.OrWhere(fmt.Sprintf("`%s` IS NULL", arg))
}

// WhereNotNull ...
func (db Database) WhereNotNull(arg string) Database {
	return db.Where(fmt.Sprintf("`%s` IS NOT NULL", arg))
}

// OrWhereNotNull ...
func (db Database) OrWhereNotNull(arg string) Database {
	return db.OrWhere(fmt.Sprintf("`%s` IS NOT NULL", arg))
}

// WhereRegexp ...
func (db Database) WhereRegexp(arg string, expstr string) Database {
	return db.Where(arg, "REGEXP", expstr)
}

// OrWhereRegexp ...
func (db Database) OrWhereRegexp(arg string, expstr string) Database {
	return db.OrWhere(arg, "REGEXP", expstr)
}

// WhereNotRegexp ...
func (db Database) WhereNotRegexp(arg string, expstr string) Database {
	return db.Where(arg, "NOT REGEXP", expstr)
}

// OrWhereNotRegexp ...
func (db Database) OrWhereNotRegexp(arg string, expstr string) Database {
	return db.OrWhere(arg, "NOT REGEXP", expstr)
}

// WhereIn ...
func (db Database) WhereIn(needle string, hystack any) Database {
	return db.Where(needle, "IN", hystack)
}

// OrWhereIn ...
func (db Database) OrWhereIn(needle string, hystack any) Database {
	return db.OrWhere(needle, "IN", hystack)
}

// WhereNotIn ...
func (db Database) WhereNotIn(needle string, hystack any) Database {
	return db.Where(needle, "NOT IN", hystack)
}

// OrWhereNotIn ...
func (db Database) OrWhereNotIn(needle string, hystack any) Database {
	return db.OrWhere(needle, "NOT IN", hystack)
}

// WhereBetween ...
func (db Database) WhereBetween(needle string, hystack []interface{}) Database {
	return db.Where(needle, "BETWEEN", hystack)
}

// OrWhereBetween ...
func (db Database) OrWhereBetween(needle string, hystack []interface{}) Database {
	return db.OrWhere(needle, "BETWEEN", hystack)
}

// WhereNotBetween ...
func (db Database) WhereNotBetween(needle string, hystack []interface{}) Database {
	return db.Where(needle, "NOT BETWEEN", hystack)
}

// OrWhereNotBetween ...
func (db Database) OrWhereNotBetween(needle string, hystack []interface{}) Database {
	return db.OrWhere(needle, "NOT BETWEEN", hystack)
}

// WhereLike ...
func (db Database) WhereLike(needle string, value string) Database {
	return db.Where(needle, "LIKE", value)
}

// OrWhereLike ...
func (db Database) OrWhereLike(needle string, value string) Database {
	return db.OrWhere(needle, "LIKE", value)
}

// WhereNotLike ...
func (db Database) WhereNotLike(needle string, value string) Database {
	return db.Where(needle, "NOT LIKE", value)
}

// OrWhereNotLike ...
func (db Database) OrWhereNotLike(needle string, value string) Database {
	return db.OrWhere(needle, "NOT LIKE", value)
}
