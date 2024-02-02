package dbgo2

func (c *Context) Table(table any, alias ...string) *Context {
	var as string
	if len(alias) > 0 {
		as = alias[0]
	}
	c.TableClause = TableClause{
		Tables: table,
		Alias:  as,
	}
	return c
}
