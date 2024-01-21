package dbgo

import (
	"gitub.com/go-webs/dbgo/iface"
)

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
func (db Database) Join(table any, argOrFn ...any) Database {
	//return db.join("INNER JOIN", table, argOrFn...)
	db.JoinBuilder.Join(table, argOrFn...)
	return db
}

// LeftJoin LEFT JOIN
func (db Database) LeftJoin(table any, argOrFn ...any) Database {
	//return db.join("LEFT JOIN", table, argOrFn...)
	db.JoinBuilder.LeftJoin(table, argOrFn...)
	return db
}

// RightJoin RIGHT JOIN
func (db Database) RightJoin(table any, argOrFn ...any) Database {
	//return db.join("RIGHT JOIN", table, argOrFn...)
	db.JoinBuilder.RightJoin(table, argOrFn...)
	return db
}

// CrossJoin CROSS JOIN
func (db Database) CrossJoin(table any) Database {
	//db.joins = append(db.joins, []interface{}{"CROSS JOIN", table})
	db.JoinBuilder.CrossJoin(table)
	return db
}

// Union expressions
func (db Database) Union(union iface.IUnion) Database {
	db.JoinBuilder.Union(union)
	return db
}
