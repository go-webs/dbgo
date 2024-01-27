package iface

type iFace interface {
	WhereClause
	Distinct() iFace
	Select(args ...any) iFace
	AddSelect(args ...any) iFace
	BuildSelect() (fields string, binds []any)
	Table(args ...any) iFace
	TableAs(arg any, as string) iFace
	Join(table interface{}, tab1Key, exp, tab2Key string) iFace
	LeftJoin(table interface{}, tab1Key, exp, tab2Key string) iFace
	RightJoin(table interface{}, tab1Key, exp, tab2Key string) iFace
	CrossJoin(table interface{}) iFace
	JoinOn(table interface{}, fn func(joins JoinClause)) iFace
	Union(iFace) (joins string)
}

type IUnion interface {
	BuildSqlQuery() (string, []any, error)
	ToSqlOnly() string
}
