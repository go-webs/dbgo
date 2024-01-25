package dbgo

// Table select from table as a
// params
//
//	arg: table string | struct bindings
//	as: table name alias
func (db Database) Table(table any, as ...string) Database {
	db.TableBuilder.Table(table, as...)
	return db
}
func (db Database) BuildTableOnly4Test() string {
	segment, _ := db.TableBuilder.BuildTable()
	return segment
}
