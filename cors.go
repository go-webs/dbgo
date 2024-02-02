package dbgo2

type IError interface {
	RecordError(error)
}

type ILogger interface {
	RecordSqlLog(query string, binds []any)
}
