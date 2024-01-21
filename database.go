package dbgo

import "gitub.com/go-webs/dbgo/iface"

type Database struct {
	*DbGo
	rawStructs
	tableStruct
	joinStruct
	selectStruct
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
