package dbgo2

type IBuilder interface {
	ToSql()
}

type IWhere interface {
	Where(column any, args ...any) IWhere
	OrWhere(column any, args ...any) IWhere
}
