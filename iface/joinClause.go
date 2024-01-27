package iface

type JoinClause interface {
	On(field1, cond, field2 string) JoinClause
	OrOn(field1, cond, field2 string) JoinClause
}
