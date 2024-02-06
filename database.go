package dbgo2

import (
	"database/sql"
	"fmt"
	"strings"
)

type Database struct {
	*DbGo
	*Session
	*Context
}

func NewDatabase(dg *DbGo) Database {
	return Database{
		dg,
		dg.NewSession(),
		&Context{Prefix: dg.prefix},
	}
}

// Distinct 在查询中添加 DISTINCT 关键字，以返回唯一结果。
func (db Database) Distinct() Database {
	db.SelectClause.Distinct = true
	return db
}

// Table sets the table name for the query.
func (db Database) Table(table any, alias ...string) Database {
	db.Context.Table(table, alias...)
	return db
}

// Select specifies the columns to retrieve.
// Select("a","b")
// Select("a.id as aid","b.id bid")
// Select("id,nickname name")
func (db Database) Select(columns ...string) Database {
	for _, column := range columns {
		splits := strings.Split(column, ",")
		for _, split := range splits {
			parts := strings.Split(strings.TrimSpace(split), " ")
			switch len(parts) {
			case 3:
				db.SelectClause.Columns = append(db.SelectClause.Columns, Column{
					Name:  strings.TrimSpace(parts[0]),
					Alias: strings.TrimSpace(parts[2]),
				})
			case 2:
				db.SelectClause.Columns = append(db.SelectClause.Columns, Column{
					Name:  strings.TrimSpace(parts[0]),
					Alias: strings.TrimSpace(parts[1]),
				})
			case 1:
				db.SelectClause.Columns = append(db.SelectClause.Columns, Column{
					Name: strings.TrimSpace(parts[0]),
				})
			}
		}
	}
	return db
}

// SelectRaw 允许直接在查询中插入原始SQL片段作为选择列。
func (db Database) SelectRaw(raw string, binds ...any) Database {
	db.SelectClause.Columns = append(db.SelectClause.Columns, Column{
		Name:  raw,
		IsRaw: true,
		Binds: binds,
	})
	return db
}

func (db Database) join(joinType string, table any, argOrFn ...any) Database {
	db.Context.JoinClause.join(joinType, table, argOrFn...)
	return db
}

// Join clause
func (db Database) Join(table any, argOrFn ...any) Database { return db.join("INNER", table, argOrFn) }
func (db Database) LeftJoin(table any, argOrFn ...any) Database {
	return db.join("LEFT", table, argOrFn)
}
func (db Database) RightJoin(table any, argOrFn ...any) Database {
	return db.join("RIGHT", table, argOrFn)
}
func (db Database) CrossJoin(table any, argOrFn ...any) Database {
	return db.join("CROSS", table, argOrFn)
}

// Where Add a basic where clause to the query.
//
//	where($column, $operator = null, $value = null, $boolean = 'and')
//
// Parameters:
//
//	array|Closure|Expression|string $column
//	mixed $operator
//	mixed $value
//	string $boolean
//
// Returns:
//
//	iface.WhereClause
//
// Examples:
//
//	Where("id=1")
//	Where("id=?",1)
//	Where("id",1)
//	Where("id","=",1)
//	Where("id","=",1,"AND")
//	Where("id","=",(select id from table limit 1))
//	Where("id","in",(select id from table), "AND")
//	Where(func(wh iface.WhereClause){wh.Where().OrWhere().WhereRaw()...})
//	Where(["id=1"])
//	Where(["id","=",1])
//	Where(["id",1])
//	Where([ ["id",1],["name","=","John"],["age",">",3] ])
func (db Database) Where(column any, argsOrclosure ...any) Database {
	db.WhereClause.Where(column, argsOrclosure...)
	return db
}
func (db Database) OrWhere(column any, argsOrclosure ...any) Database {
	db.WhereClause.OrWhere(column, argsOrclosure...)
	return db
}

// WhereRaw 在查询中添加一个原生SQL“where”条件。
//
// sql: 原生SQL条件字符串。
// bindings: SQL绑定参数数组。
func (db Database) WhereRaw(raw string, bindings ...any) Database {
	db.WhereClause.WhereRaw(raw, bindings...)
	return db
}
func (db Database) OrWhereRaw(raw string, bindings ...any) Database {
	db.WhereClause.OrWhereRaw(raw, bindings...)
	return db
}

// GroupBy 添加 GROUP BY 子句
func (db Database) GroupBy(columns ...string) Database {
	db.Groups = append(db.Groups, columns...)
	return db
}

// Having 添加 HAVING 子句, 同where
func (db Database) Having(column any, argsOrClosure ...any) Database {
	db.HavingClause.Where(column, argsOrClosure...)
	return db
}
func (db Database) OrHaving(column any, argsOrClosure ...any) Database {
	db.HavingClause.OrWhere(column, argsOrClosure...)
	return db
}

// HavingRaw 添加 HAVING 子句, 同where
func (db Database) HavingRaw(raw string, argsOrClosure ...any) Database {
	db.HavingClause.WhereRaw(raw, argsOrClosure...)
	return db
}
func (db Database) OrHavingRaw(raw string, argsOrClosure ...any) Database {
	db.HavingClause.OrWhereRaw(raw, argsOrClosure...)
	return db
}

// OrderBy adds an ORDER BY clause to the query.
func (db Database) OrderBy(column string, directions ...string) Database {
	var direction string
	if len(directions) > 0 {
		direction = directions[0]
	}
	db.OrderByClause.Columns = append(db.OrderByClause.Columns, OrderByItem{
		Column:    column,
		Direction: direction,
	})
	return db
}
func (db Database) OrderByRaw(column string) Database {
	db.OrderByClause.Columns = append(db.OrderByClause.Columns, OrderByItem{
		Column: column,
		IsRaw:  true,
	})
	return db
}

// Limit 设置查询结果的限制数量。
func (db Database) Limit(limit int) Database {
	db.LimitOffsetClause.Limit = limit
	return db
}

// Offset 设置查询结果的偏移量。
func (db Database) Offset(offset int) Database {
	db.LimitOffsetClause.Offset = offset
	return db
}

// Page 页数,根据limit确定
func (db Database) Page(num int) Database {
	db.LimitOffsetClause.Page = num
	return db
}

// Get 获取查询结果集。
//
// columns: 要获取的列名数组，如果不提供，则获取所有列。
func (db Database) Get(columns ...string) (res []map[string]any, err error) {
	var prepare string
	var binds []any
	prepare, binds, err = db.Select(columns...).ToSql()
	if err != nil {
		return
	}

	err = db.queryToBindResult(&res, prepare, binds...)
	return
}
func (db Database) First(columns ...string) (res map[string]any, err error) {
	var prepare string
	var binds []any
	prepare, binds, err = db.Select(columns...).Limit(1).ToSql()
	if err != nil {
		return
	}

	res = make(map[string]any)
	err = db.queryToBindResult(&res, prepare, binds...)
	return
}
func (db Database) Find(id int) (res map[string]any, err error) {
	var prepare string
	var binds []any
	prepare, binds, err = db.Where("id", id).Limit(1).ToSql()
	if err != nil {
		return
	}

	res = make(map[string]any)
	err = db.queryToBindResult(&res, prepare, binds...)
	return
}
func (db Database) To(obj any, mustFields ...string) (err error) {
	var prepare string
	var binds []any
	prepare, binds, err = db.ToSqlTo(obj, mustFields...)
	if err != nil {
		return
	}

	err = db.queryToBindResult(obj, prepare, binds...)
	return
}
func (db Database) queryToBindResult(bind any, query string, args ...any) (err error) {
	return db.Session.QueryTo(bind, query, args...)
}
func (db Database) insert(obj any, ignoreCase string, onDuplicateKeys []string, mustFields ...string) (res sql.Result, err error) {
	segment, binds, err := db.ToSqlInsert(obj, ignoreCase, onDuplicateKeys, mustFields...)
	if err != nil {
		return res, err
	}
	return db.Session.Exec(segment, binds...)
}
func (db Database) InsertGetId(obj any, mustFields ...string) (lastInsertId int64, err error) {
	result, err := db.insert(obj, "", []string{}, mustFields...)
	if err != nil {
		return lastInsertId, err
	}
	return result.LastInsertId()
}
func (db Database) Insert(obj any, mustFields ...string) (aff int64, err error) {
	result, err := db.insert(obj, "", []string{}, mustFields...)
	if err != nil {
		return aff, err
	}
	return result.RowsAffected()
}
func (db Database) InsertOrIgnore(obj any, mustFields ...string) (aff int64, err error) {
	result, err := db.insert(obj, "IGNORE", []string{}, mustFields...)
	if err != nil {
		return aff, err
	}
	return result.RowsAffected()
}
func (db Database) Upsert(obj any, onDuplicateKeys []string, mustFields ...string) (aff int64, err error) {
	result, err := db.insert(obj, "IGNORE", onDuplicateKeys, mustFields...)
	if err != nil {
		return aff, err
	}
	return result.RowsAffected()
}
func (db Database) UpdateOrInsert(attributes, values map[string]any) (affectedRows int64, err error) {
	dbTmp := db.Where(attributes)
	var exists bool
	if exists, err = dbTmp.Exists(); err != nil {
		return
	}
	if exists {
		return dbTmp.Update(values)
	}
	return dbTmp.Insert(values)
}
func (db Database) Update(obj any, mustFields ...string) (aff int64, err error) {
	segment, binds, err := db.ToSqlUpdate(obj, mustFields...)
	if err != nil {
		return aff, err
	}
	return db.Session.execute(segment, binds...)
}
func (db Database) Delete(obj any) (aff int64, err error) {
	segment, binds, err := db.ToSqlDelete(obj)
	if err != nil {
		return aff, err
	}
	return db.Session.execute(segment, binds...)
}
func (db Database) incDecEach(symbol string, data map[string]any) (aff int64, err error) {
	prepare, values, err := db.ToSqlIncDec(symbol, data)
	if err != nil {
		return aff, err
	}
	return db.Session.execute(prepare, values...)
}
func (db Database) incDec(symbol string, column string, steps ...any) (aff int64, err error) {
	var step any = 1
	if len(steps) > 0 {
		step = steps[0]
	}
	return db.incDecEach(symbol, map[string]any{column: step})
}
func (db Database) Increment(column string, steps ...any) (aff int64, err error) {
	return db.incDec("+", column, steps...)
}
func (db Database) Decrement(column string, steps ...any) (aff int64, err error) {
	return db.incDec("-", column, steps...)
}
func (db Database) IncrementEach(data map[string]any) (aff int64, err error) {
	return db.incDecEach("+", data)
}
func (db Database) DecrementEach(data map[string]any) (aff int64, err error) {
	return db.incDecEach("-", data)
}

// func (db Database) Aggregate(functions, columns string) (float64, error) {}
func (db Database) aggregateSingle(bind any, function, column string) error {
	prepare, values, err := db.ToSqlAggregate(function, column)
	if err != nil {
		return err
	}
	return db.Session.QueryRow(prepare, values...).Scan(bind)
}
func (db Database) Max(column string) (res float64, err error) {
	err = db.aggregateSingle(&res, "max", column)
	return
}
func (db Database) Min(column string) (res float64, err error) {
	err = db.aggregateSingle(&res, "min", column)
	return
}
func (db Database) Sum(column string) (res float64, err error) {
	err = db.aggregateSingle(&res, "sum", column)
	return
}
func (db Database) Avg(column string) (res float64, err error) {
	err = db.aggregateSingle(&res, "avg", column)
	return
}
func (db Database) Count() (res int64, err error) {
	err = db.aggregateSingle(&res, "count", "*")
	return
}

// Pluck 从查询结果集中获取指定列的值列表。
func (db Database) Pluck(column string) (res []any, err error) {
	ress, err := db.Get(column)
	if err != nil {
		return res, err
	}
	for _, v := range ress {
		res = append(res, v[column])
	}
	return
}

// List 获取指定列的键值对列表。
func (db Database) List(column string, keyColumn string) (res []map[any]any, err error) {
	ress, err := db.Get(column, keyColumn)
	if err != nil {
		return res, err
	}
	for _, v := range ress {
		res = append(res, map[any]any{
			v[keyColumn]: v[column],
		})
	}
	return
}
func (db Database) Value(column string) (res any, err error) {
	first, err := db.First(column)
	if err != nil {
		return res, err
	}
	return first[column], err
}
func (db Database) Exists(bind ...any) (b bool, err error) {
	prepare, values, err := db.ToSqlExists(bind...)
	if err != nil {
		return b, err
	}
	err = db.Session.QueryRow(prepare, values...).Scan(&b)
	return
}
func (db Database) DoesntExist(bind ...any) (b bool, err error) {
	b, err = db.Exists(bind...)
	return !b, err
}
func (db Database) Union(b IBuilder, unionAll ...bool) (res []map[string]any, err error) {
	prepare, values, err := db.ToSql()
	if err != nil {
		return res, err
	}
	sql4prepare, binds, err := b.ToSql()
	if err != nil {
		return res, err
	}
	var union = "UNION"
	if len(unionAll) > 0 && unionAll[0] {
		union = "UNION ALL"
	}
	err = db.queryToBindResult(&res, fmt.Sprintf("%s %s %s", prepare, union, sql4prepare), append(values, binds...))
	return
}
func (db Database) UnionAll(b IBuilder) (res []map[string]any, err error) {
	return db.Union(b, true)
}

func (db Database) Truncate(obj ...any) (affectedRows int64, err error) {
	var table string
	var dbTmp = db
	if len(obj) > 0 {
		dbTmp = db.Table(obj[0])
	}
	table, _, err = dbTmp.ToSqlTable()
	if err != nil {
		return
	}
	return db.Session.execute(fmt.Sprintf("TRUNCATE TABLE %s", BackQuotes(table)))
}
