package dbgo2

type Context struct {
	SelectClause
	TableClause
	JoinClause
	WhereClause
	HavingClause
	OrderByClause
	LimitOffsetClause
	Groups []string

	Prefix  string
	Queries string
	Args    []any
	Err     error
}
