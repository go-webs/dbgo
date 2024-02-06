package dbgo2

import "log/slog"

func LogCors() HandlerFunc {
	return func(c *Context) {
		slog.Debug("sql:%s, %v\n", c.Queries, c.Bindings)
	}
}

func ErrorCors() HandlerFunc {
	return func(c *Context) {
		slog.Error("error:%s", c.Err)
	}
}

func DeleteCors() HandlerFunc {
	return func(c *Context) {
		c.Where("is_deleted", 0)
	}
}
