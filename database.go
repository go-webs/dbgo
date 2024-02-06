package dbgo2

import (
	"database/sql"
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

// Join clause
func (db Database) Join(table any, argOrFn ...any) Database {
	db.Context.JoinClause.join("INNER", table, argOrFn...)
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

// WhereRaw 在查询中添加一个原生SQL“where”条件。
//
// sql: 原生SQL条件字符串。
// bindings: SQL绑定参数数组。
func (db Database) WhereRaw(raw string, bindings ...any) Database {
	db.WhereClause.WhereRaw(raw, bindings...)
	return db
}

// GroupBy 添加 GROUP BY 子句
func (db Database) GroupBy(columns ...string) Database {
	db.Groups = append(db.Groups, columns...)
	return db
}

// Having 添加 HAVING 子句, 同where
func (db Database) Having(column any, argsOrclosure ...any) Database {
	db.HavingClause.Where(column, argsOrclosure...)
	return db
}

// HavingRaw 添加 HAVING 子句, 同where
func (db Database) HavingRaw(raw string, argsOrclosure ...any) Database {
	db.HavingClause.WhereRaw(raw, argsOrclosure...)
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
func (db Database) insert(obj any, mustFields ...string) (res sql.Result, err error) {
	segment, binds, err := db.ToSqlInsert(obj, mustFields...)
	if err != nil {
		return res, err
	}
	return db.Session.Exec(segment, binds...)
}
func (db Database) Insert(obj any, mustFields ...string) (aff int64, err error) {
	result, err := db.insert(obj, mustFields...)
	if err != nil {
		return aff, err
	}
	return result.RowsAffected()
}
func (db Database) InsertGetId(obj any, mustFields ...string) (lastInsertId int64, err error) {
	result, err := db.insert(obj, mustFields...)
	if err != nil {
		return lastInsertId, err
	}
	return result.LastInsertId()
}
func (db Database) Update(obj any, mustFields ...string) (aff int64, err error) {
	segment, binds, err := db.ToSqlUpdate(obj, mustFields...)
	if err != nil {
		return aff, err
	}
	result, err := db.Session.Exec(segment, binds...)
	if err != nil {
		return aff, err
	}
	return result.RowsAffected()
}
func (db Database) Delete(obj any) (aff int64, err error) {
	segment, binds, err := db.ToSqlDelete(obj)
	if err != nil {
		return aff, err
	}
	result, err := db.Session.Exec(segment, binds...)
	if err != nil {
		return aff, err
	}
	return result.RowsAffected()
}
