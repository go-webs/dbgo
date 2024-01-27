package dbgo

import (
	"fmt"
	"github.com/go-webs/dbgo/util"
)

func (db Database) GroupBy(args ...string) Database {
	db.GroupBuilder.GroupBy(args...)
	return db
}
func (db Database) Having(args ...any) Database {
	db.GroupBuilder.Having(args...)
	return db
}

// GroupByRaw fields
// params
//
//	arg: expressions
func (db Database) GroupByRaw(args ...string) Database {
	db.GroupBuilder.GroupByRaw(args...)
	return db
}

// HavingRaw fields
// params
//
//	arg: expressions
//	binds: bind values
func (db Database) HavingRaw(arg string, binds ...any) Database {
	db.GroupBuilder.HavingRaw(arg, binds)
	return db
}

func (db Database) HavingBetween(col string, haystack ...any) Database {
	db.HavingRaw(fmt.Sprintf("%s BETWEEN ? AND ?", util.BackQuotes(col)), haystack)
	return db
}
func (db Database) HavingNotBetween(col string, haystack ...any) Database {
	db.HavingRaw(fmt.Sprintf("%s NOT BETWEEN ? AND ?", util.BackQuotes(col)), haystack)
	return db
}
func (db Database) HavingNull(col string) Database {
	db.HavingRaw(fmt.Sprintf("%s IS NULL", util.BackQuotes(col)))
	return db
}
func (db Database) HavingNotNull(col string) Database {
	db.HavingRaw(fmt.Sprintf("%s IS NOT NULL", util.BackQuotes(col)))
	return db
}

//func (db Database) HavingNested(col string, haystack ...any) Database {
//	return db
//}

func (db Database) OrHaving(args ...any) Database {
	db.GroupBuilder.OrHaving(args...)
	return db
}

func (db Database) OrHavingRaw(args ...any) Database {
	db.GroupBuilder.OrHavingRaw(args...)
	return db
}

func (db Database) OrHavingBetween(col string, haystack ...any) Database {
	db.OrHavingRaw(fmt.Sprintf("%s BETWEEN ? AND ?", util.BackQuotes(col)), haystack)
	return db
}
func (db Database) OrHavingNotBetween(col string, haystack ...any) Database {
	db.OrHavingRaw(fmt.Sprintf("%s NOT BETWEEN ? AND ?", util.BackQuotes(col)), haystack)
	return db
}
func (db Database) OrHavingNull(col string) Database {
	db.OrHavingRaw(fmt.Sprintf("%s IS NULL", util.BackQuotes(col)))
	return db
}
func (db Database) OrHavingNotNull(col string) Database {
	db.OrHavingRaw(fmt.Sprintf("%s IS NOT NULL", util.BackQuotes(col)))
	return db
}
