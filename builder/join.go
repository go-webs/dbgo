package builder

import (
	"errors"
	"fmt"
	"github.com/go-webs/dbgo/iface"
	"github.com/go-webs/dbgo/util"
	"strings"
)

type JoinBuilder struct {
	joins      []joinStruct
	joinClause [][]any
	union      iface.IUnion
	prefix     string
}
type joinStruct struct {
	joinType string
	table    *TableBuilder
	argOrFn  []any
}

func NewJoinBuilder(prefix string) *JoinBuilder {
	return &JoinBuilder{prefix: prefix}
}

func (jc *JoinBuilder) On(tab1Key, cond, tab2Key string) iface.JoinClause {
	jc.joinClause = append(jc.joinClause, []interface{}{"AND", tab1Key, cond, tab2Key})
	return jc
}

func (jc *JoinBuilder) OrOn(tab1Key, cond, tab2Key string) iface.JoinClause {
	jc.joinClause = append(jc.joinClause, []interface{}{"OR", tab1Key, cond, tab2Key})
	return jc
}

func (jc *JoinBuilder) join(joinType string, table any, argOrFn ...any) *JoinBuilder {
	if v, ok := table.(*TableBuilder); ok {
		v.prefix = jc.prefix
		jc.joins = append(jc.joins, joinStruct{joinType: joinType, table: v, argOrFn: argOrFn})
	} else {
		jc.joins = append(jc.joins, joinStruct{joinType: joinType, table: NewTableBuilder(jc.prefix).Table(table), argOrFn: argOrFn})
	}
	return jc
}

// Join INNER JOIN
// param
//
//	table: string | TableAs()=>[]interface{}
//	tab1Key: a.id
//	exp: =
//	tab2Key: b.aid
//
// example
//
//	db.Table("user").Join("card", "user.id","=","card.user_id")
//	db.Table("user","a").Join(TableAs("card","b"), "a.id","=","b.user_id")
//	db.Table(User{},"a").Join(TableAs(Card{},"b"), "a.id","=","b.user_id")
//	db.Table(User{},"a").Join(TableAs(Card{},"b"), func(dbgo.JohnClause))
func (jc *JoinBuilder) Join(table any, argOrFn ...any) *JoinBuilder {
	if len(argOrFn) == 1 {
		return jc.join("JOIN", table, argOrFn...)
	}
	return jc.join("INNER JOIN", table, argOrFn...)
}

// LeftJoin LEFT JOIN
func (jc *JoinBuilder) LeftJoin(table any, argOrFn ...any) *JoinBuilder {
	return jc.join("LEFT JOIN", table, argOrFn...)
}

// RightJoin RIGHT JOIN
func (jc *JoinBuilder) RightJoin(table any, argOrFn ...any) *JoinBuilder {
	return jc.join("RIGHT JOIN", table, argOrFn...)
}

// CrossJoin CROSS JOIN
func (jc *JoinBuilder) CrossJoin(table any) *JoinBuilder {
	return jc.join("CROSS JOIN", table)
}

//func (jc *JoinBuilder) JoinOn(table any, fn func(joins iface.JoinClause)) *JoinBuilder {
//	return jc.join("JOIN", table, fn)
//}

// BuildJoin expressions
func (jc *JoinBuilder) BuildJoin() (joins string, bindValues []any, err error) {
	if jc.union != nil {
		joins, bindValues, err = jc.union.BuildSqlQuery()
		if err != nil {
			return
		}
		//joins = jc.union.ToSqlOnly()
		joins = fmt.Sprintf("UNION (%s)", joins)
		return
	}
	for _, v := range jc.joins {
		var tab string
		tab, err = v.table.BuildTable()
		if err != nil {
			return
		}
		if v.joinType == "CROSS JOIN" {
			joins = fmt.Sprintf("%s %s %s", joins, v.joinType, tab)
		} else if v.joinType == "JOIN" {
			jc.joinClause = [][]any{}
			if joinTmp, ok := v.argOrFn[0].(func(joins iface.JoinClause)); ok {
				joinTmp(jc)
			} else {
				err = errors.New("the second param of join must be func in two params")
				return
			}
			var tmp []string
			for _, v2 := range jc.joinClause {
				if len(tmp) == 0 { // 第一个不加 and 或者 or 链接
					tmp = append(tmp, fmt.Sprintf("%s %s %s", util.BackQuotes(v2[1]), v2[2], util.BackQuotes(v2[3])))
				} else {
					tmp = append(tmp, fmt.Sprintf("%s %s %s %s", v2[0], util.BackQuotes(v2[1]), v2[2], util.BackQuotes(v2[3])))
				}
			}
			joins = fmt.Sprintf("%s %s %s ON (%s)", joins, v.joinType, tab, strings.Join(tmp, " "))
		} else {
			if len(v.argOrFn) == 3 {
				joins = fmt.Sprintf("%s %s %s ON %s %s %s", joins, v.joinType, tab, util.BackQuotes(v.argOrFn[0]), v.argOrFn[1], util.BackQuotes(v.argOrFn[2]))
			} else {
				err = errors.New("the param's length after first must be 3")
				return
			}
		}
	}
	return strings.TrimSpace(joins), bindValues, err
}

// Union expressions
func (jc *JoinBuilder) Union(union iface.IUnion) *JoinBuilder {
	jc.union = union
	return jc
}
