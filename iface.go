package dbgo

type IWhere interface {
	Where(column any, args ...any) IWhere
	OrWhere(column any, args ...any) IWhere
	WhereRaw(raw string, bindingsAndBoolean ...any) IWhere
	OrWhereRaw(sqlSeg string, bindingsAndBoolean ...any) IWhere

	WhereBetween(column string, values any) IWhere
	OrWhereBetween(column string, values any) IWhere
	WhereNotBetween(column string, values any) IWhere
	OrWhereNotBetween(column string, values any) IWhere

	WhereIn(column string, values any) IWhere
	OrWhereIn(column string, values any) IWhere
	WhereNotIn(column string, values any) IWhere
	OrWhereNotIn(column string, values any) IWhere

	WhereNull(column string) IWhere
	OrWhereNull(column string) IWhere
	WhereNotNull(column string) IWhere
	OrWhereNotNull(column string) IWhere

	WhereLike(column string, value string) IWhere
	OrWhereLike(column string, value string) IWhere
	WhereNotLike(column string, value string) IWhere
	OrWhereNotLike(column string, value string) IWhere

	WhereExists(clause IBuilder) IWhere
	WhereNotExists(clause IBuilder) IWhere
}

type IJoin interface {
	On(column string, args ...string) IJoin
	OrOn(column string, args ...string) IJoin
}

type IBuilder interface {
	//Distinct() Database
	//Tables(table any, alias ...string) Database
	//Select(columns ...string) Database
	//SelectRaw(raw string, binds ...any) Database
	//Where(column any, argsOrclosure ...any) Database
	//WhereRaw(raw string, bindings ...any) Database
	//GroupBy(columns ...string) Database
	//Having(column any, argsOrclosure ...any) Database
	//HavingRaw(raw string, argsOrclosure ...any) Database
	//OrderBy(column string, directions ...string) Database
	//Limit(limit int) Database
	//Offset(offset int) Database
	//Page(num int) Database
	//Get(columns ...string) (res []map[string]any, Err error)
	//First() (res map[string]any, Err error)

	ToSql() (sql4prepare string, binds []any, err error)
	ToSqlIncDec(symbol string, data map[string]any) (sql4prepare string, values []any, err error)
	ToSqlSelect() (sql4prepare string, binds []any)
	ToSqlTable() (sql4prepare string, values []any, err error)
	ToSqlJoin() (sql4prepare string, binds []any, err error)
	ToSqlWhere() (sql4prepare string, values []any, err error)
	ToSqlOrderBy() (sql4prepare string)
	ToSqlLimitOffset() (sqlSegment string, binds []any)
	ToSqlInsert(obj any, ignore string, onDuplicateKeys []string, mustFields ...string) (sqlSegment string, binds []any, err error)
	ToSqlUpdate(obj any, mustFields ...string) (sqlSegment string, binds []any, err error)
	ToSqlDelete(obj any) (sqlSegment string, binds []any, err error)
}

type IDriver interface {
	ToSql(c *Context) (sql4prepare string, binds []any, err error)
	ToSqlIncDec(c *Context, symbol string, data map[string]any) (sql4prepare string, values []any, err error)
	ToSqlSelect(c *Context) (sql4prepare string, binds []any)
	ToSqlTable(c *Context) (sql4prepare string, values []any, err error)
	ToSqlJoin(c *Context) (sql4prepare string, binds []any, err error)
	ToSqlWhere(c *Context) (sql4prepare string, values []any, err error)
	ToSqlOrderBy(c *Context) (sql4prepare string)
	ToSqlLimitOffset(c *Context) (sqlSegment string, binds []any)
	ToSqlInsert(c *Context, obj any, ignore string, onDuplicateKeys []string, mustFields ...string) (sqlSegment string, binds []any, err error)
	ToSqlUpdate(c *Context, obj any, mustFields ...string) (sqlSegment string, binds []any, err error)
	ToSqlDelete(c *Context, obj any) (sqlSegment string, binds []any, err error)
}
