package mysql

import (
	"errors"
	"fmt"
	"go-webs/dbgo2"
	"reflect"
	"strings"
)

const DriverName = "mysql"

type Builder struct {
	prefix string
}

func init() {
	dbgo2.Register(DriverName, &Builder{})
}

func (b Builder) ToSql(c *dbgo2.Context) (sql4prepare string, binds []any, err error) {
	selects, anies := b.ToSqlSelect(c)
	table, binds2, err := b.ToSqlTable(c)
	if err != nil {
		return sql4prepare, binds2, err
	}
	where, binds3, err := b.ToSqlWhere(c)
	if err != nil {
		return sql4prepare, binds2, err
	}
	orderBy := b.ToSqlOrderBy(c)
	limit, binds4 := b.ToSqlLimitOffset(c)

	binds = append(binds, anies...)
	binds = append(binds, binds3...)
	binds = append(binds, binds4...)
	sql4prepare = fmt.Sprintf("%s FROM %s %s %s %s", selects, table, where, orderBy, limit)
	return
}

func (Builder) ToSqlSelect(c *dbgo2.Context) (sql4prepare string, binds []any) {
	var cols []string
	for _, col := range c.SelectClause.Columns {
		if col.IsRaw {
			cols = append(cols, col.Name)
			binds = append(binds, col.Binds...)
		} else {
			if col.Alias == "" {
				cols = append(cols, BackQuotes(col.Name))
			} else {
				cols = append(cols, fmt.Sprintf("%s AS %s", BackQuotes(col.Name), col.Alias))
			}
		}
	}
	var distinct string
	if c.SelectClause.Distinct {
		distinct = "DISTINCT "
	}
	sql4prepare = fmt.Sprintf("SELECT %s%s", distinct, strings.Join(cols, ", "))
	return
}

func (b Builder) ToSqlTable(c *dbgo2.Context) (sql4prepare string, binds []any, err error) {
	if v, ok := c.TableClause.Tables.(dbgo2.IBuilder); ok {
		return v.ToSql()
	}
	rfv := reflect.Indirect(reflect.ValueOf(c.TableClause.Tables))
	switch rfv.Kind() {
	case reflect.String:
		sql4prepare = BackQuotes(fmt.Sprintf("%s%s", c.Prefix, c.TableClause.Tables))
	case reflect.Struct:
		sql4prepare = b.buildTableName(rfv.Type(), c.Prefix)
	case reflect.Slice:
		if rfv.Type().Elem().Kind() == reflect.Struct {
			sql4prepare = b.buildTableName(rfv.Type().Elem(), c.Prefix)
		} else {
			c.Err = errors.New("table param must be string or struct(slice) bind with 1 or 2 params")
			return
		}
	default:
		c.Err = errors.New("table must string | struct | slice")
		return
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", sql4prepare, c.TableClause.Alias)), binds, err
}

func (b Builder) ToSqlWhere(c *dbgo2.Context) (sql4prepare string, binds []any, err error) {
	if len(c.WhereClause.Conditions) == 0 {
		return
	}
	var sql4prepareArr []string
	for _, v := range c.Conditions {
		switch v.(type) {
		case dbgo2.TypeWhereRaw:
			item := v.(dbgo2.TypeWhereRaw)
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s", item.LogicalOp, item.Column))
			binds = append(binds, item.Bindings...)
		case dbgo2.TypeWhereNested:
			item := v.(dbgo2.TypeWhereNested)
			var tmp = dbgo2.Context{}
			item.Column(&tmp.WhereClause)
			prepare, anies, err := b.ToSqlWhere(&tmp)
			if err != nil {
				return sql4prepare, binds, err
			}
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s (%s)", item.LogicalOp, prepare))
			binds = append(binds, anies...)
		case dbgo2.TypeWhereSubQuery:
			item := v.(dbgo2.TypeWhereSubQuery)
			query, anies, err := item.SubQuery.ToSql()
			if err != nil {
				return sql4prepare, binds, err
			}
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s (%s)", item.LogicalOp, BackQuotes(item.Column), item.Operator, query))
			binds = append(binds, anies...)
		case dbgo2.TypeWhereStandard:
			item := v.(dbgo2.TypeWhereStandard)
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s ?", item.LogicalOp, BackQuotes(item.Column), item.Operator))
			binds = append(binds, item.Value)
		case dbgo2.TypeWhereIn:
			item := v.(dbgo2.TypeWhereIn)
			values := ToSlice(item.Value)
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s (%s)", item.LogicalOp, BackQuotes(item.Column), item.Operator, strings.Repeat("?,", len(values)-1)+"?"))
			binds = append(binds, values...)
		case dbgo2.TypeWhereBetween:
			item := v.(dbgo2.TypeWhereBetween)
			values := ToSlice(item.Value)
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s ? AND ?", item.LogicalOp, BackQuotes(item.Column), item.Operator))
			binds = append(binds, values...)
		}
	}
	if len(sql4prepareArr) > 0 {
		whereTmp := strings.TrimSpace(strings.Trim(strings.Trim(strings.TrimSpace(strings.Join(sql4prepareArr, " ")), "AND"), "OR"))
		sql4prepare = fmt.Sprintf("WHERE %s", whereTmp)
	}
	return
}

func (b Builder) ToSqlOrderBy(c *dbgo2.Context) (sql4prepare string) {
	if len(c.OrderByClause.Columns) == 0 {
		return
	}
	var orderBys []string
	for _, v := range c.OrderByClause.Columns {
		if v.Direction == "" {
			orderBys = append(orderBys, BackQuotes(v.Column))
		} else {
			orderBys = append(orderBys, fmt.Sprintf("%s %s", BackQuotes(v.Column), v.Direction))
		}
	}
	sql4prepare = fmt.Sprintf("ORDER BY %s", strings.Join(orderBys, ", "))
	return
}

func (b Builder) ToSqlLimitOffset(c *dbgo2.Context) (sqlSegment string, binds []any) {
	var offset int
	if c.LimitOffsetClause.Offset > 0 {
		offset = c.LimitOffsetClause.Offset
	} else if c.LimitOffsetClause.Page > 0 {
		offset = c.LimitOffsetClause.Limit * (c.LimitOffsetClause.Page - 1)
	}
	if c.LimitOffsetClause.Limit > 0 {
		if offset > 0 {
			sqlSegment = "LIMIT ? OFFSET ?"
			binds = append(binds, c.LimitOffsetClause.Limit, offset)
		} else {
			sqlSegment = "LIMIT ?"
			binds = append(binds, c.LimitOffsetClause.Limit)
		}
	}
	return
}

func (b Builder) ToSqlInsert(c *dbgo2.Context, obj any, mustFields ...string) (sqlSegment string, binds []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Struct:
		var datas []map[string]any
		datas, err = dbgo2.StructsToInsert(obj, mustFields...)
		if err != nil {
			return
		}
		c.Table(obj)
		return b.toSqlInsert(c, datas, "")
	case reflect.Slice:
		switch rfv.Type().Elem().Kind() {
		case reflect.Struct:
			c.Table(obj)
			var datas []map[string]any
			datas, err = dbgo2.StructsToInsert(obj, mustFields...)
			if err != nil {
				return
			}
			return b.toSqlInsert(c, datas, "")
		default:
			return b.toSqlInsert(c, obj, "")
		}
	default:
		return b.toSqlInsert(c, obj, "")
	}
}

func (b Builder) ToSqlUpdate(c *dbgo2.Context, obj any, mustFields ...string) (sqlSegment string, binds []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Struct:
		dataMap, pk, pkValue, err := dbgo2.StructToUpdate(obj, mustFields...)
		if err != nil {
			return sqlSegment, binds, err
		}
		c.Table(obj)
		if pk != "" {
			c.Where(pk, pkValue)
		}
		return b.toSqlUpdate(c, dataMap)
	default:
		return b.toSqlUpdate(c, obj)
	}
}

func (b Builder) ToSqlDelete(c *dbgo2.Context, obj any) (sqlSegment string, binds []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Struct:
		pk, pkValue := dbgo2.StructToDelete(obj)
		if err != nil {
			return sqlSegment, binds, err
		}
		c.Table(obj)
		if pk != "" {
			c.Where(pk, pkValue)
		}
		return b.toSqlSqlDelete(c)
	default:
		return b.toSqlSqlDelete(c)
	}
}
