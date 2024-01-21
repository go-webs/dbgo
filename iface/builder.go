package iface

type Builder interface {
	Where(args ...interface{}) Builder
	OrWhere(args ...interface{}) Builder
	// null
	WhereNull(arg string) Builder
	OrWhereNull(arg string) Builder
	WhereNotNull(arg string) Builder
	OrWhereNotNull(arg string) Builder
	// like
	WhereLike(needle string, value string) Builder
	OrWhereLike(needle string, value string) Builder
	WhereNotLike(needle string, value string) Builder
	OrWhereNotLike(needle string, value string) Builder
	// regexp
	WhereRegexp(arg string, expstr string) Builder
	OrWhereRegexp(arg string, expstr string) Builder
	WhereNotRegexp(arg string, expstr string) Builder
	OrWhereNotRegexp(arg string, expstr string) Builder
	// in
	WhereIn(needle string, hystack []interface{}) Builder
	OrWhereIn(needle string, hystack []interface{}) Builder
	WhereNotIn(needle string, hystack []interface{}) Builder
	OrWhereNotIn(needle string, hystack []interface{}) Builder
	// between
	WhereBetween(needle string, hystack []interface{}) Builder
	OrWhereBetween(needle string, hystack []interface{}) Builder
	WhereNotBetween(needle string, hystack []interface{}) Builder
	OrWhereNotBetween(needle string, hystack []interface{}) Builder
	// exists
	WhereExists(where Builder) Builder

	Distinct() Builder
	Select(args ...any) Builder
	AddSelect(args ...any) Builder
	BuildSelect() (fields string, binds []any)
	Table(args ...any) Builder
	TableAs(arg any, as string) Builder
	Join(table interface{}, tab1Key, exp, tab2Key string) Builder
	LeftJoin(table interface{}, tab1Key, exp, tab2Key string) Builder
	RightJoin(table interface{}, tab1Key, exp, tab2Key string) Builder
	CrossJoin(table interface{}) Builder
	JoinOn(table interface{}, fn func(joins JoinClause)) Builder
	Union(Builder) (joins string)
}
