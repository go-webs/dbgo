package dbgo

type rawStruct struct {
	expression string
	binds      []any
}
type rawStructs struct {
	selectRaw  []rawStruct
	whereRaw   []rawStruct
	havingRaw  []rawStruct
	orderByRaw []string
	groupByRaw []string
}

// SelectRaw fields
// params
//
//	arg: expressions
//	binds: bind values
func (db Database) SelectRaw(arg string, binds ...any) Database {
	db.selectRaw = append(db.selectRaw, rawStruct{arg, binds})
	return db
}

// WhereRaw fields
// params
//
//	arg: expressions
//	binds: bind values
func (db Database) WhereRaw(arg string, binds ...any) Database {
	db.whereRaw = append(db.whereRaw, rawStruct{arg, binds})
	return db
}

// HavingRaw fields
// params
//
//	arg: expressions
//	binds: bind values
func (db Database) HavingRaw(arg string, binds ...any) Database {
	db.havingRaw = append(db.havingRaw, rawStruct{arg, binds})
	return db
}

// OrderByRaw fields
// params
//
//	arg: expressions
func (db Database) OrderByRaw(arg string) Database {
	db.orderByRaw = append(db.orderByRaw, arg)
	return db
}

// GroupByRaw fields
// params
//
//	arg: expressions
func (db Database) GroupByRaw(arg string) Database {
	db.groupByRaw = append(db.groupByRaw, arg)
	return db
}
