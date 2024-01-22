package iface

type WhereClause interface {
	Where(args ...interface{}) WhereClause
	OrWhere(args ...interface{}) WhereClause
	WhereRaw(arg string, binds ...any) WhereClause
	OrWhereRaw(arg string, binds ...any) WhereClause
	//// null
	//WhereNull(arg string) WhereClause
	//OrWhereNull(arg string) WhereClause
	//WhereNotNull(arg string) WhereClause
	//OrWhereNotNull(arg string) WhereClause
	//// like
	//WhereLike(needle string, value string) WhereClause
	//OrWhereLike(needle string, value string) WhereClause
	//WhereNotLike(needle string, value string) WhereClause
	//OrWhereNotLike(needle string, value string) WhereClause
	//// regexp
	//WhereRegexp(arg string, expstr string) WhereClause
	//OrWhereRegexp(arg string, expstr string) WhereClause
	//WhereNotRegexp(arg string, expstr string) WhereClause
	//OrWhereNotRegexp(arg string, expstr string) WhereClause
	//// in
	//WhereIn(needle string, hystack []interface{}) WhereClause
	//OrWhereIn(needle string, hystack []interface{}) WhereClause
	//WhereNotIn(needle string, hystack []interface{}) WhereClause
	//OrWhereNotIn(needle string, hystack []interface{}) WhereClause
	//// between
	//WhereBetween(needle string, hystack []interface{}) WhereClause
	//OrWhereBetween(needle string, hystack []interface{}) WhereClause
	//WhereNotBetween(needle string, hystack []interface{}) WhereClause
	//OrWhereNotBetween(needle string, hystack []interface{}) WhereClause
	//// exists
	//WhereExists(where WhereClause) WhereClause
}
