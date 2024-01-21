package dbgo

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

var operator = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "not like",
	"in", "not in", "between", "not between", "regexp", "not regexp"}

// Where : query or execute where condition, the relation is and
func (db Database) Where(args ...interface{}) Database {
	if len(args) == 0 {
		return db
	}
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组
	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	w := []interface{}{"AND", args}
	db.where = append(db.where, w)
	return db
}

// Where : query or execute where condition, the relation is and
func (db Database) OrWhere(args ...interface{}) Database {
	if len(args) == 0 {
		return db
	}
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组
	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	w := []interface{}{"OR", args}

	db.where = append(db.where, w)

	return db
}

// WhereNull ...
func (db Database) WhereNull(arg string) Database {
	return db.Where(arg + " IS NULL")
}

// OrWhereNull ...
func (db Database) OrWhereNull(arg string) Database {
	return db.OrWhere(arg + " IS NULL")
}

// WhereNotNull ...
func (db Database) WhereNotNull(arg string) Database {
	return db.Where(arg + " IS NOT NULL")
}

// OrWhereNotNull ...
func (db Database) OrWhereNotNull(arg string) Database {
	return db.OrWhere(arg + " IS NOT NULL")
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
func (db Database) WhereIn(needle string, hystack []interface{}) Database {
	return db.Where(needle, "IN", hystack)
}

// OrWhereIn ...
func (db Database) OrWhereIn(needle string, hystack []interface{}) Database {
	return db.OrWhere(needle, "IN", hystack)
}

// WhereNotIn ...
func (db Database) WhereNotIn(needle string, hystack []interface{}) Database {
	return db.Where(needle, "NOT IN", hystack)
}

// OrWhereNotIn ...
func (db Database) OrWhereNotIn(needle string, hystack []interface{}) Database {
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

// parseWhere : parse where condition
func (db Database) parseWhere() (string, error) {
	// 取出所有where
	wheres := db.where
	// where解析后存放每一项的容器
	var where []string

	for _, args := range wheres {
		// and或者or条件
		var condition = args[0].(string)
		// 统计当前数组中有多少个参数
		params := ToSlice(args[1])
		paramsLength := len(params)

		switch paramsLength {
		case 3: // 常规3个参数:  {"id",">",1}
			res, err := db.parseParams(params)
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)

		case 2: // 常规2个参数:  {"id",1}
			res, err := db.parseParams(params)
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
					whereArr = append(whereArr, db.AddFieldQuotes(key)+"="+db.GetPlaceholder())
					db.SetBindValues(val)
				}
				if len(whereArr) != 0 {
					where = append(where, condition+" ("+strings.Join(whereArr, " and ")+")")
				}
			case []interface{}: // 一维数组
				var whereArr []string
				whereMoreLength := len(paramReal)
				switch whereMoreLength {
				case 3, 2, 1:
					res, err := db.parseParams(paramReal)
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
						res, err := db.parseParams(arr)
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
			case func(Database):
				// 清空where,给嵌套的where让路,复用这个节点
				db.where = [][]interface{}{}

				// 执行嵌套where放入Database struct
				paramReal(db)
				// 再解析一遍后来嵌套进去的where
				whereNested, err := db.parseWhere()
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
func (db Database) parseParams(args []interface{}) (s string, err error) {
	paramsLength := len(args)
	argsReal := args

	// 存储当前所有数据的数组
	var paramsToArr []string

	switch paramsLength {
	case 3: // 常规3个参数:  {"id",">",1}
		//if !inArray(argsReal[1], b.GetRegex()) {
		if !SliceContains(db.GetOperator(), strings.ToLower(argsReal[1].(string))) {
			err = errors.New("where parameter is wrong")
			return
		}

		//paramsToArr = append(paramsToArr, argsReal[0].(string))
		paramsToArr = append(paramsToArr, db.AddFieldQuotes(argsReal[0].(string)))
		paramsToArr = append(paramsToArr, argsReal[1].(string))

		switch strings.Trim(strings.ToLower(reflect.ValueOf(argsReal[1]).String()), " ") {
		//case "like", "not like":
		//	paramsToArr = append(paramsToArr, b.GetPlaceholder())
		//	b.SetBindValues(argsReal[2])
		case "in", "not in":
			var tmp []string
			//reflect.Indirect(reflect.ValueOf(argsReal[2]))
			var ar2 = ToSlice(argsReal[2])
			for _, item := range ar2 {
				tmp = append(tmp, db.GetPlaceholder())
				db.SetBindValues(item)
			}
			paramsToArr = append(paramsToArr, "("+strings.Join(tmp, ",")+")")

		case "between", "not between":
			var ar2 = ToSlice(argsReal[2])
			paramsToArr = append(paramsToArr, db.GetPlaceholder()+" AND "+db.GetPlaceholder())
			db.SetBindValues(ar2[0])
			db.SetBindValues(ar2[1])

		default:
			paramsToArr = append(paramsToArr, db.GetPlaceholder())
			db.SetBindValues(argsReal[2])
		}
	case 2:
		paramsToArr = append(paramsToArr, db.AddFieldQuotes(argsReal[0].(string)))
		paramsToArr = append(paramsToArr, "=")
		paramsToArr = append(paramsToArr, db.GetPlaceholder())
		db.SetBindValues(argsReal[1])
	case 1:
		paramsToArr = append(paramsToArr, argsReal[0].(string))
	}

	return strings.Join(paramsToArr, " "), nil
}

func (db Database) BuildWhereOnly() (where string, values []interface{}, err error) {
	// 存储原values
	var valuesClone = slices.Clone(db.whereBindValues)

	where, err = db.parseWhere()
	if err != nil {
		return
	}
	for _, v := range db.whereBindValues {
		values = append(values, v)
	}

	// 还原values
	db.whereBindValues = slices.Clone(valuesClone)
	return
}

// AddFieldQuotes ...
func (db Database) AddFieldQuotes(field string) string {
	reg := regexp.MustCompile(`^\w+$`)
	if reg.MatchString(field) {
		return fmt.Sprintf("`%s`", field)
	}
	return field
}

// GetOperator ...
func (db Database) GetOperator() []string {
	return operator
}

// GetPlaceholder ...
func (db Database) GetPlaceholder() string {
	return "?"
}

// SetBindValues ...
func (db Database) SetBindValues(arg interface{}) {
	db.whereBindValues = append(db.whereBindValues, arg)
}
