package dbgo2

type IError interface {
	RecordError(error)
}

type ILog interface {
	RecordSqlLog(query string, binds []any)
}
