package iface

type TypeWhere interface {
	string | []interface{} | [][]interface{}
}
type Builder interface {
	Where(args ...interface{}) Builder
	OrWhere(args ...interface{}) Builder
}
type JoinClause interface {
	On(field1, cond, field2 string) JoinClause
	OrOn(field1, cond, field2 string) JoinClause
}
