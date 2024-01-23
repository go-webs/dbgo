package dbgo

import (
	"database/sql"
	"gitub.com/go-webs/dbgo/builder"
)

type transaction struct {
	tx *sql.Tx
}

func (db *transaction) Begin(sd *sql.DB) (err error) {
	db.tx, err = sd.Begin()
	return
}

// Rollback ...
func (db *transaction) Rollback() (err error) {
	err = db.tx.Rollback()
	db.tx = nil
	return
}

// Commit ...
func (db *transaction) Commit() (err error) {
	err = db.tx.Commit()
	db.tx = nil
	return
}

type Database struct {
	*DbGo
	rawStructs
	//tableStruct
	//joinStruct
	//selectStruct
	//selects  []any
	distinct string
	//tables          tableStruct
	//where           [][]interface{}
	//whereBindValues []interface{}
	//union           iface.IUnion

	*builder.TableBuilder
	*builder.SelectBuilder
	*builder.JoinBuilder
	*builder.WhereBuilder
	*builder.GroupBuilder
	*builder.OrderByBuilder
	*builder.PageBuilder

	sharedLock    string
	lockForUpdate string
	*transaction
}

func NewDB(dg *DbGo) *Database {
	return &Database{
		DbGo:           dg,
		TableBuilder:   builder.NewTableBuilder(dg.Cluster.Prefix),
		SelectBuilder:  builder.NewSelectBuilder(),
		JoinBuilder:    builder.NewJoinBuilder(dg.Cluster.Prefix),
		WhereBuilder:   builder.NewWhereBuilder(),
		GroupBuilder:   builder.NewGroupBuilder(),
		OrderByBuilder: builder.NewOrderByBuilder(),
		PageBuilder:    builder.NewPageBuilder(),
	}
}

// Distinct for distinct
func (db Database) Distinct() Database {
	db.distinct = "DISTINCT"
	return db
}

// SharedLock 4 lock in share mode
func (db Database) SharedLock() Database {
	db.sharedLock = "LOCK IN SHARE MODE"
	return db
}

// LockForUpdate 4 for update
func (db Database) LockForUpdate() Database {
	db.lockForUpdate = "FOR UPDATE"
	return db
}

// Begin ...
func (db Database) Begin() (err error) {
	err = db.transaction.Begin(db.getMasterDB())
	return
}
