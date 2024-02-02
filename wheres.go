package dbgo2

import (
	"errors"
	"reflect"
	"slices"
	"strings"
)

//	type WhereClause struct {
//		Conditions []any
//	}
func (w *WhereClause) addTypeWhereRaw(boolean string, value string, bindings []any) *WhereClause {
	w.Conditions = append(w.Conditions, TypeWhereRaw{LogicalOp: boolean, Column: value, Bindings: bindings})
	return w
}
func (w *WhereClause) addTypeWhereNested(boolean string, value func(where IWhere)) *WhereClause {
	w.Conditions = append(w.Conditions, TypeWhereNested{LogicalOp: boolean, Column: value})
	return w
}
func (w *WhereClause) addTypeWhereSubQuery(boolean string, column string, operator string, value IBuilder) *WhereClause {
	w.Conditions = append(w.Conditions, TypeWhereSubQuery{LogicalOp: boolean, Column: column, Operator: operator, SubQuery: value})
	return w
}
func (w *WhereClause) addTypeWhereIn(boolean string, column string, operator string, value []any) *WhereClause {
	w.Conditions = append(w.Conditions, TypeWhereIn{LogicalOp: boolean, Column: column, Operator: operator, Value: value})
	return w
}
func (w *WhereClause) addTypeWhereBetween(boolean string, column string, operator string, value []any) *WhereClause {
	w.Conditions = append(w.Conditions, TypeWhereBetween{LogicalOp: boolean, Column: column, Operator: operator, Value: value})
	return w
}
func (w *WhereClause) addTypeWhereStandard(boolean string, column string, operator string, value any) *WhereClause {
	w.Conditions = append(w.Conditions, TypeWhereStandard{LogicalOp: boolean, Column: column, Operator: operator, Value: value})
	return w
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
//	SubQuery
//
// Laravel api
func (w *WhereClause) WhereRaw(raw string, bindingsAndBoolean ...any) IWhere {
	return w.whereRaw("AND", raw, bindingsAndBoolean...)
}

// OrWhereRaw clause
func (w *WhereClause) OrWhereRaw(sqlSeg string, bindingsAndBoolean ...any) IWhere {
	return w.whereRaw("OR", sqlSeg, bindingsAndBoolean...)
}

func (w *WhereClause) whereRaw(boolean string, sqlSeg string, bindingsAndBoolean ...any) IWhere {
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
func (w *WhereClause) Where(column any, args ...any) IWhere {
	return w.where("AND", column, args...)
}

// OrWhere clause
func (w *WhereClause) OrWhere(column any, args ...any) IWhere {
	return w.where("OR", column, args...)
}

func (w *WhereClause) where(boolean string, column any, args ...any) IWhere {
	if column == nil {
		return w
	}
	switch len(args) {
	case 0:
		rfv := reflect.Indirect(reflect.ValueOf(column))
		switch rfv.Kind() {
		case reflect.Func:
			if fn, ok := column.(func(where IWhere)); ok {
				w.addTypeWhereNested(boolean, fn)
			} else {
				w.Err = errors.New("not supported where params")
			}
		case reflect.String:
			return w.whereRaw(boolean, rfv.String())
		case reflect.Slice:
			if rfv.Len() > 1 {
				rfvItem := rfv.Index(0)
				if rfvItem.Kind() == reflect.Slice {
					return w.where(boolean, rfvItem.Interface())
				} else {
					var tmp []any
					for i := 0; i < rfv.Len(); i++ {
						tmp = append(tmp, rfv.Index(i).Interface())
					}
					return w.where(boolean, tmp[0], tmp[1:]...)
				}
			} else if rfv.Len() > 0 {
				return w.whereRaw(boolean, rfv.Index(0).String())
			}
			w.Err = errors.New("not supported where params")
		default:
			w.Err = errors.New("not supported where params")
		}
	case 1:
		if IsExpression(column) {
			return w.whereRaw(boolean, column.(string), args...)
		}
		return w.where(boolean, column, "=", args[0], boolean)
	case 2:
		return w.where(boolean, column, args[0], args[1], boolean)
	case 3:
		rfv := reflect.Indirect(reflect.ValueOf(args[1]))
		if rfv.Kind() == reflect.Slice { // in/between
			var operators = []string{"in", "not in"}
			if slices.Contains(operators, strings.ToLower(args[0].(string))) {
				val := ToSlice(args[1])
				if len(val) > 0 {
					w.addTypeWhereIn(args[2].(string), column.(string), args[0].(string), ToSlice(args[1]))
				}
			}
			operators = []string{"between", "not between"}
			if slices.Contains(operators, strings.ToLower(args[0].(string))) {
				val := ToSlice(args[1])
				if len(val) > 0 {
					w.addTypeWhereBetween(args[2].(string), column.(string), args[0].(string), ToSlice(args[1]))
				}
			}
		} else if builder, ok := args[1].(IBuilder); ok {
			w.addTypeWhereSubQuery(args[2].(string), column.(string), args[0].(string), builder)
		} else {
			w.addTypeWhereStandard(args[2].(string), column.(string), args[0].(string), args[1])
		}
	default:
		w.Err = errors.New("not supported where params")
	}
	return w
}
