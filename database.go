package dbgo

import (
	"gitub.com/go-webs/dbgo/builder"
)

type Database struct {
	*DbGo
	rawStructs
	//tableStruct
	//joinStruct
	//selectStruct
	//selects  []any
	distinct string
	//tables          tableStruct
	where           [][]interface{}
	whereBindValues []interface{}
	//union           iface.IUnion

	*builder.TableBuilder
	*builder.SelectBuilder
	*builder.JoinBuilder
	*builder.WhereBuilder
}

func NewDB(dg *DbGo) *Database {
	return &Database{
		DbGo:          dg,
		TableBuilder:  builder.NewTableBuilder(dg.Cluster.Prefix),
		SelectBuilder: builder.NewSelectBuilder(),
		JoinBuilder:   builder.NewJoinBuilder(dg.Cluster.Prefix),
		WhereBuilder:  builder.NewWhereBuilder(),
	}
}

// Distinct for distinct
func (db Database) Distinct() Database {
	db.distinct = "DISTINCT"
	return db
}
