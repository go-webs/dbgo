package builder

import (
	"errors"
	"fmt"
	"github.com/go-webs/dbgo/iface"
	"github.com/go-webs/dbgo/util"
	"reflect"
	"slices"
	"strings"
)

type typeWhereRaw struct {
	boolean  string
	column   string
	bindings []any
}
type typeWhereNested struct {
	boolean string
	column  func(iface.WhereClause2)
}
type typeWhereSubQuery struct {
	boolean  string
	column   any
	operator string
	value    iface.IUnion
}
type typeWhereStandard struct {
	boolean  string
	column   any
	operator string
	value    any
}
type typeWhereIn struct {
	boolean  string
	column   any
	operator string
	value    []any
}
type typeWhereBetween struct {
	boolean  string
	column   any
	operator string
	value    []any
}

// WhereBuilderNew struct
type WhereBuilderNew struct {
	wheres []any
	err    error
	Exists bool
}

// NewWhereBuilderNew ptr
func NewWhereBuilderNew() *WhereBuilderNew {
	return &WhereBuilderNew{}
}

func (w *WhereBuilderNew) addTypeWhereRaw(boolean string, value string, bindings []any) *WhereBuilderNew {
	w.wheres = append(w.wheres, typeWhereRaw{boolean: boolean, column: value, bindings: bindings})
	return w
}

func (w *WhereBuilderNew) addTypeWhereNested(boolean string, value func(iface.WhereClause2)) *WhereBuilderNew {
	w.wheres = append(w.wheres, typeWhereNested{boolean: boolean, column: value})
	return w
}
func (w *WhereBuilderNew) addTypeWhereSubQuery(boolean string, column any, operator string, value iface.IUnion) *WhereBuilderNew {
	w.wheres = append(w.wheres, typeWhereSubQuery{boolean: boolean, column: column, operator: operator, value: value})
	return w
}
func (w *WhereBuilderNew) addTypeWhereIn(boolean string, column any, operator string, value []any) *WhereBuilderNew {
	w.wheres = append(w.wheres, typeWhereIn{boolean: boolean, column: column, operator: operator, value: value})
	return w
}
func (w *WhereBuilderNew) addTypeWhereBetween(boolean string, column any, operator string, value []any) *WhereBuilderNew {
	w.wheres = append(w.wheres, typeWhereBetween{boolean: boolean, column: column, operator: operator, value: value})
	return w
}
func (w *WhereBuilderNew) addTypeWhereStandard(boolean string, column any, operator string, value any) *WhereBuilderNew {
	w.wheres = append(w.wheres, typeWhereStandard{boolean: boolean, column: column, operator: operator, value: value})
	return w
}

func (w *WhereBuilderNew) BuildWhere() (sql4prepare string, values []any, err error) {
	if w.err != nil {
		return sql4prepare, values, w.err
	}
	if len(w.wheres) > 0 {
		w.Exists = true
	}

	var sql4prepareArr []string
	for _, v := range w.wheres {
		switch v.(type) {
		case typeWhereRaw:
			item := v.(typeWhereRaw)
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s", item.boolean, item.column))
			values = append(values, item.bindings...)
		case typeWhereNested:
			item := v.(typeWhereNested)
			var tmp = NewWhereBuilderNew()
			item.column(tmp)
			prepare, anies, err := tmp.BuildWhere()
			if err != nil {
				return sql4prepare, values, err
			}
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s (%s)", item.boolean, prepare))
			values = append(values, anies...)
		case typeWhereSubQuery:
			item := v.(typeWhereSubQuery)
			query, anies, err := item.value.BuildSqlQuery()
			if err != nil {
				return sql4prepare, values, err
			}
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s (%s)", item.boolean, util.BackQuotes(item.column), item.operator, query))
			values = append(values, anies...)
		case typeWhereStandard:
			item := v.(typeWhereStandard)
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s ?", item.boolean, util.BackQuotes(item.column), item.operator))
			values = append(values, item.value)
		case typeWhereIn:
			item := v.(typeWhereIn)
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s (%s)", item.boolean, util.BackQuotes(item.column), item.operator, strings.Repeat("?,", len(item.value)-1)+"?"))
			values = append(values, item.value...)
		case typeWhereBetween:
			item := v.(typeWhereBetween)
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s ? AND ?", item.boolean, util.BackQuotes(item.column), item.operator))
			values = append(values, item.value...)
		}
	}
	sql4prepare = strings.TrimSpace(strings.Trim(strings.Trim(strings.TrimSpace(strings.Join(sql4prepareArr, " ")), "AND"), "OR"))
	return
}

// WhereRaw Add a raw where clause to the query.
//
//	whereRaw($sql, $bindings = [], $boolean = 'and')
//
// Parameters:
//
//	string $sql
//	mixed $bindings
//	string $boolean
//
// Returns:
//
//	Builder
//
// Laravel api
func (w *WhereBuilderNew) WhereRaw(sqlSeg string, bindingsAndBoolean ...any) iface.WhereClause2 {
	return w.whereRaw("AND", sqlSeg, bindingsAndBoolean...)
}

// OrWhereRaw clause
func (w *WhereBuilderNew) OrWhereRaw(sqlSeg string, bindingsAndBoolean ...any) iface.WhereClause2 {
	return w.whereRaw("OR", sqlSeg, bindingsAndBoolean...)
}
func (w *WhereBuilderNew) whereRaw(boolean string, sqlSeg string, bindingsAndBoolean ...any) iface.WhereClause2 {
	if sqlSeg == "" {
		return w
	}
	if len(bindingsAndBoolean) == 0 {
		return w.WhereRaw(sqlSeg, []any{}, boolean)
	} else if len(bindingsAndBoolean) == 1 {
		return w.WhereRaw(sqlSeg, bindingsAndBoolean[0], boolean)
	} else if len(bindingsAndBoolean) == 2 {
		rfv := reflect.ValueOf(bindingsAndBoolean[0])
		var bindTmp []any
		if rfv.Kind() == reflect.Slice {
			for i := 0; i < rfv.Len(); i++ {
				bindTmp = append(bindTmp, rfv.Index(i).Interface())
			}
		} else {
			bindTmp = append(bindTmp, rfv.Interface())
		}
		rfv1 := reflect.ValueOf(bindingsAndBoolean[1])
		if rfv1.Kind() == reflect.String {
			w.addTypeWhereRaw(rfv1.String(), sqlSeg, bindTmp)
		}
	}
	return w
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
func (w *WhereBuilderNew) Where(column any, args ...any) iface.WhereClause2 {
	return w.where("AND", column, args...)
}

// OrWhere clause
func (w *WhereBuilderNew) OrWhere(column any, args ...any) iface.WhereClause2 {
	return w.where("OR", column, args...)
}
func (w *WhereBuilderNew) where(boolean string, column any, args ...any) iface.WhereClause2 {
	if column == nil {
		return w
	}
	switch len(args) {
	case 0:
		rfv := reflect.Indirect(reflect.ValueOf(column))
		switch rfv.Kind() {
		case reflect.Func:
			if fn, ok := column.(func(iface.WhereClause2)); ok {
				w.addTypeWhereNested(boolean, fn)
			} else {
				w.err = errors.New("not supported where params")
			}
		case reflect.String:
			return w.WhereRaw(rfv.String())
		case reflect.Slice:
			if rfv.Len() > 1 {
				rfvItem := rfv.Index(0)
				if rfvItem.Kind() == reflect.Slice {
					for i := 0; i < rfv.Len(); i++ {
						w.Where(rfv.Index(i).Interface())
					}
				} else {
					var tmp []any
					for i := 0; i < rfv.Len(); i++ {
						tmp = append(tmp, rfv.Index(i).Interface())
					}
					w.Where(tmp[0], tmp[1:]...)
				}
			} else if rfv.Len() > 0 {
				return w.WhereRaw(rfv.Index(0).String())
			}
			w.err = errors.New("not supported where params")
		default:
			w.err = errors.New("not supported where params")
			return w
		}
	case 1:
		if iface.IsExpression(column) {
			return w.whereRaw(boolean, column.(string), args...)
		}
		return w.Where(column, "=", args[0], boolean)
	case 2:
		return w.Where(column, args[0], args[1], boolean)
	case 3:
		rfv := reflect.Indirect(reflect.ValueOf(args[1]))
		if rfv.Kind() == reflect.Slice { // in/between
			var operators = []string{"in", "not in"}
			if slices.Contains(operators, strings.ToLower(args[0].(string))) {
				val := util.ToSlice(args[1])
				if len(val) > 0 {
					w.addTypeWhereIn(args[2].(string), column, args[0].(string), util.ToSlice(args[1]))
				}
			}
			operators = []string{"between", "not between"}
			if slices.Contains(operators, strings.ToLower(args[0].(string))) {
				val := util.ToSlice(args[1])
				if len(val) > 0 {
					w.addTypeWhereBetween(args[2].(string), column, args[0].(string), util.ToSlice(args[1]))
				}
			}
		} else if builder, ok := args[1].(iface.IUnion); ok {
			w.addTypeWhereSubQuery(args[2].(string), column, args[0].(string), builder)
		} else {
			w.addTypeWhereStandard(args[2].(string), column, args[0].(string), args[1])
		}
	default:
		w.err = errors.New("not supported where params")
	}
	return w
}
