package dbgo2

type IWhere interface {
	Where(column any, args ...any) IWhere
	OrWhere(column any, args ...any) IWhere
	WhereRaw(raw string, bindingsAndBoolean ...any) IWhere
	OrWhereRaw(sqlSeg string, bindingsAndBoolean ...any) IWhere
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
	//Where(column any, argsOrCloser ...any) Database
	//WhereRaw(raw string, bindings ...any) Database
	//GroupBy(columns ...string) Database
	//Having(column any, argsOrCloser ...any) Database
	//HavingRaw(raw string, argsOrCloser ...any) Database
	//OrderBy(column string, directions ...string) Database
	//Limit(limit int) Database
	//Offset(offset int) Database
	//Page(num int) Database
	//Get(columns ...string) (res []map[string]any, Err error)
	//First() (res map[string]any, Err error)

	ToSql() (sql4prepare string, binds []any, err error)
	ToSqlSelect() (sql4prepare string, binds []any)
	ToSqlTable() (sql4prepare string, values []any, err error)
	ToSqlWhere() (sql4prepare string, values []any, err error)
	ToSqlOrderBy() (sql4prepare string)
	ToSqlLimitOffset() (sqlSegment string, binds []any)
	ToSqlInsert(obj any, mustFields ...string) (sqlSegment string, binds []any, err error)
	ToSqlUpdate(obj any, mustFields ...string) (sqlSegment string, binds []any, err error)
	ToSqlDelete(obj any) (sqlSegment string, binds []any, err error)
}

type IDriver interface {
	ToSql(ctx *Context) (sql4prepare string, binds []any, err error)
	ToSqlSelect(c *Context) (sql4prepare string, binds []any)
	ToSqlTable(w *Context) (sql4prepare string, values []any, err error)
	ToSqlWhere(w *Context) (sql4prepare string, values []any, err error)
	ToSqlOrderBy(c *Context) (sql4prepare string)
	ToSqlLimitOffset(c *Context) (sqlSegment string, binds []any)
	ToSqlInsert(c *Context, obj any, mustFields ...string) (sqlSegment string, binds []any, err error)
	ToSqlUpdate(c *Context, obj any, mustFields ...string) (sqlSegment string, binds []any, err error)
	ToSqlDelete(c *Context, obj any) (sqlSegment string, binds []any, err error)
}
