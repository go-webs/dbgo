package dbgo

import (
	"fmt"
	"gitub.com/go-webs/dbgo/iface"
	"reflect"
	"strings"
)

type joinStruct struct {
	joins      [][]any
	joinClause [][]any
}

func (j *joinStruct) On(tab1Key, cond, tab2Key string) iface.JoinClause {
	j.joinClause = append(j.joinClause, []interface{}{"AND", tab1Key, cond, tab2Key})
	return j
}

func (j *joinStruct) OrOn(tab1Key, cond, tab2Key string) iface.JoinClause {
	j.joinClause = append(j.joinClause, []interface{}{"OR", tab1Key, cond, tab2Key})
	return j
}

func (db Database) join(joinType string, table any, tab1Key, exp, tab2Key string) Database {
	if reflect.Indirect(reflect.ValueOf(table)).Kind() == reflect.String {
		table = []interface{}{table}
	}
	db.joins = append(db.joins, []any{joinType, table, tab1Key, exp, tab2Key})
	return db
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
func (db Database) Join(table any, tab1Key, exp, tab2Key string) Database {
	return db.join("INNER JOIN", table, tab1Key, exp, tab2Key)
}

// LeftJoin LEFT JOIN
func (db Database) LeftJoin(table any, tab1Key, exp, tab2Key string) Database {
	return db.join("LEFT JOIN", table, tab1Key, exp, tab2Key)
}

// RightJoin RIGHT JOIN
func (db Database) RightJoin(table any, tab1Key, exp, tab2Key string) Database {
	return db.join("RIGHT JOIN", table, tab1Key, exp, tab2Key)
}

// CrossJoin CROSS JOIN
func (db Database) CrossJoin(table any) Database {
	db.joins = append(db.joins, []interface{}{"CROSS JOIN", table})
	return db
}
func (db Database) JoinOn(table any, fn func(joins iface.JoinClause)) Database {
	db.joins = append(db.joins, []interface{}{"JOIN", table, fn})
	return db
}

// BuildJoin expressions
func (db Database) BuildJoin() (joins string) {
	if db.union != nil {
		return fmt.Sprintf("UNION ALL (%s)", db.union.ToSql())
	}
	for _, v := range db.joins {
		tab := newTable(db.DbGo.Cluster.Prefix, v[1]).buildTable()
		if v[0] == "CROSS JOIN" {
			joins = fmt.Sprintf("%s %s %s", joins, v[0], tab)
		} else if v[0] == "JOIN" {
			db.joinStruct.joinClause = [][]interface{}{}
			v[2].(func(joins iface.JoinClause))(&db.joinStruct)
			var tmp []string
			for _, v2 := range db.joinClause {
				if len(tmp) == 0 { // 第一个不加 and 或者 or 链接
					tmp = append(tmp, fmt.Sprintf("%s %s %s", v2[1], v2[2], v2[3]))
				} else {
					tmp = append(tmp, fmt.Sprintf("%s %s %s %s", v2[0], v2[1], v2[2], v2[3]))
				}
			}
			joins = fmt.Sprintf("%s %s %s ON (%s)", joins, v[0], tab, strings.Join(tmp, " "))
		} else {
			joins = fmt.Sprintf("%s %s %s ON %s %s %s", joins, v[0], tab, v[2], v[3], v[4])
		}
	}
	return strings.TrimSpace(joins)
}

// Union expressions
func (db Database) Union(union iface.IUnion) Database {
	db.union = union
	return db
}
