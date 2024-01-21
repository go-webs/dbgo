package dbgo

import "gitub.com/go-webs/dbgo/iface"

type Database struct {
	*DbGo
	rawStruct
	tableStruct
	joinStruct
	selects  []any
	distinct string
	//tables          tableStruct
	where           [][]interface{}
	whereBindValues []interface{}
	union           iface.IUnion
}

// Distinct for distinct
func (db Database) Distinct() Database {
	db.distinct = "DISTINCT"
	return db
}
