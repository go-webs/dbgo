package builder

import (
	"errors"
	"fmt"
	"gitub.com/go-webs/dbgo/iface"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

var operator = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "not like",
	"in", "not in", "between", "not between", "regexp", "not regexp"}

type whereStruct struct {
	relation    string
	expressions []any
}
type whereRawStruct struct {
	relation    string
	expressions rawStruct
}
type WhereBuilder struct {
	wheres          []whereStruct
	wheresRaw       []whereRawStruct
	whereBindValues []any
}

func NewWhereBuilder() *WhereBuilder {
	return &WhereBuilder{}
}

func (w *WhereBuilder) Where(args ...any) iface.WhereClause {
	w.wheres = append(w.wheres, whereStruct{
		relation:    "AND",
		expressions: args,
	})
	return w
}

func (w *WhereBuilder) OrWhere(args ...any) iface.WhereClause {
	w.wheres = append(w.wheres, whereStruct{
		relation:    "OR",
		expressions: args,
	})
	return w
}

func (w *WhereBuilder) WhereRaw(arg string, binds ...any) iface.WhereClause {
	w.wheresRaw = append(w.wheresRaw, whereRawStruct{
		relation: "AND",
		expressions: rawStruct{
			expression: arg,
			binds:      binds,
		},
	})
	return w
}

func (w *WhereBuilder) OrWhereRaw(arg string, binds ...any) iface.WhereClause {
	w.wheresRaw = append(w.wheresRaw, whereRawStruct{
		relation: "OR",
		expressions: rawStruct{
			expression: arg,
			binds:      binds,
		},
	})
	return w
}

// parseWhere : parse where condition
func (w *WhereBuilder) parseWhere() (string, error) {
	// 取出所有where
	wheres := w.wheres
	// where解析后存放每一项的容器
	var where []string

	for _, args := range wheres {
		// and或者or条件
		var condition = args.relation
		// 统计当前数组中有多少个参数
		params := ToSlice(args.expressions)
		paramsLength := len(params)

		switch paramsLength {
		case 3: // 常规3个参数:  {"id",">",1}  或者  子查询
			res, err := w.parseParams(params)
			if err != nil {
				return res, err
			}
			where = append(where, fmt.Sprintf("%s %s", condition, res))

		case 2: // 常规2个参数:  {"id",1}
			res, err := w.parseParams([]any{params[0], "=", params[1]})
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)
		case 1: // 二维数组或字符串
			switch paramReal := params[0].(type) {
			case string:
				where = append(where, condition+" ("+paramReal+")")
			case map[string]interface{}: // map
				var whereArr []string
				for key, val := range paramReal {
					whereArr = append(whereArr, w.AddFieldQuotes(key)+"="+w.GetPlaceholder())
					w.SetBindValues(val)
				}
				if len(whereArr) != 0 {
					where = append(where, condition+" ("+strings.Join(whereArr, " and ")+")")
				}
			case []interface{}: // 一维数组
				var whereArr []string
				whereMoreLength := len(paramReal)
				switch whereMoreLength {
				case 3, 2, 1:
					res, err := w.parseParams(paramReal)
					if err != nil {
						return res, err
					}
					whereArr = append(whereArr, res)
				default:
					return "", errors.New("where data format is wrong")
				}
				if len(whereArr) != 0 {
					where = append(where, condition+" ("+strings.Join(whereArr, " and ")+")")
				}
			case [][]interface{}: // 二维数组
				var whereMore []string
				for _, arr := range paramReal { // {{"a", 1}, {"id", ">", 1}}
					whereMoreLength := len(arr)
					switch whereMoreLength {
					case 3, 2, 1:
						res, err := w.parseParams(arr)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					default:
						return "", errors.New("where data format is wrong")
					}
				}
				if len(whereMore) != 0 {
					where = append(where, condition+" ("+strings.Join(whereMore, " and ")+")")
				}
			case func(wh iface.WhereClause):
				// 清空where,给嵌套的where让路,复用这个节点
				w.wheres = []whereStruct{}

				// 执行嵌套where放入WhereBuilder struct
				paramReal(w)
				// 再解析一遍后来嵌套进去的where
				whereNested, err := w.parseWhere()
				if err != nil {
					return "", err
				}
				// 嵌套的where放入一个括号内
				where = append(where, condition+" ("+whereNested+")")
			default:
				return "", errors.New("where data format is wrong")
			}
		}
	}

	// 合并where,去掉左侧的空格,and,or并返回
	return strings.TrimLeft(
		strings.TrimPrefix(
			strings.TrimPrefix(
				strings.Trim(
					strings.Join(where, " "),
					" "),
				"AND"),
			"OR"),
		" "), nil
}

/**
 * 将where条件中的参数转换为where条件字符串
 * example: {"id",">",1}, {"age", 18}
 */
// parseParams : 将where条件中的参数转换为where条件字符串
func (w *WhereBuilder) parseParams(args []any) (s string, err error) {
	paramsLength := len(args)
	argsReal := args

	// 存储当前所有数据的数组
	var paramsToArr []string

	switch paramsLength {
	case 3: // 常规3个参数:  {"id",">",1}
		//if !inArray(argsReal[1], b.GetRegex()) {
		if !slices.Contains(w.GetOperator(), strings.ToLower(argsReal[1].(string))) {
			err = errors.New("where parameter is wrong")
			return
		}

		//paramsToArr = append(paramsToArr, argsReal[0].(string))
		paramsToArr = append(paramsToArr, w.AddFieldQuotes(argsReal[0].(string)))
		paramsToArr = append(paramsToArr, argsReal[1].(string))

		switch strings.Trim(strings.ToLower(reflect.ValueOf(argsReal[1]).String()), " ") {
		//case "like", "not like":
		//	paramsToArr = append(paramsToArr, b.GetPlaceholder())
		//	b.SetBindValues(argsReal[2])
		case "in", "not in":
			switch reflect.TypeOf(argsReal[2]).Kind() {
			case reflect.Slice:
				var tmp []string
				var ar2 = ToSlice(argsReal[2])
				for _, item := range ar2 {
					tmp = append(tmp, w.GetPlaceholder())
					w.SetBindValues(item)
				}
				paramsToArr = append(paramsToArr, "("+strings.Join(tmp, ",")+")")
			default: // sub query
				if v, ok := argsReal[2].(iface.IUnion); ok {
					query, anies, err2 := v.BuildSqlQuery()
					if err2 != nil {
						return s, err2
					}
					paramsToArr = append(paramsToArr, fmt.Sprintf("(%s)", query))
					w.SetBindValues(anies...)
				}
			}

		case "between", "not between":
			var ar2 = ToSlice(argsReal[2])
			paramsToArr = append(paramsToArr, w.GetPlaceholder()+" AND "+w.GetPlaceholder())
			w.SetBindValues(ar2[0])
			w.SetBindValues(ar2[1])

		default:
			// sub query
			if v, ok := argsReal[2].(iface.IUnion); ok {
				query, anies, err2 := v.BuildSqlQuery()
				if err2 != nil {
					return s, err2
				}
				paramsToArr = append(paramsToArr, fmt.Sprintf("(%s)", query))
				w.SetBindValues(anies...)
			} else {
				paramsToArr = append(paramsToArr, w.GetPlaceholder())
				w.SetBindValues(argsReal[2])
			}
		}
	case 2:
		paramsToArr = append(paramsToArr, w.AddFieldQuotes(argsReal[0].(string)))
		paramsToArr = append(paramsToArr, "=")
		paramsToArr = append(paramsToArr, w.GetPlaceholder())
		w.SetBindValues(argsReal[1])
	case 1:
		paramsToArr = append(paramsToArr, argsReal[0].(string))
	}

	return strings.Join(paramsToArr, " "), nil
}

func (w *WhereBuilder) BuildWhereOnly() (where string) {
	// 存储原values
	var valuesClone = slices.Clone(w.whereBindValues)

	where, _, _ = w.BuildWhere()

	// 还原values
	w.whereBindValues = slices.Clone(valuesClone)
	return
}
func (w *WhereBuilder) BuildWhere() (where string, values []interface{}, err error) {
	where, err = w.parseWhere()
	if err != nil {
		return
	}
	for _, v := range w.whereBindValues {
		values = append(values, v)
	}
	return
}

// AddFieldQuotes ...
func (w *WhereBuilder) AddFieldQuotes(field string) string {
	reg := regexp.MustCompile(`^\w+$`)
	if reg.MatchString(field) {
		return fmt.Sprintf("`%s`", field)
	}
	return field
}

// GetOperator ...
func (w *WhereBuilder) GetOperator() []string {
	return operator
}

// GetPlaceholder ...
func (w *WhereBuilder) GetPlaceholder() string {
	return "?"
}

// SetBindValues ...
func (w *WhereBuilder) SetBindValues(args ...any) {
	w.whereBindValues = append(w.whereBindValues, args...)
}
