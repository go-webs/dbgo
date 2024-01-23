package builder

import (
	"fmt"
	"gitub.com/go-webs/dbgo/util"
	"strings"
)

type OrderByBuilder struct {
	args []string
}

func NewOrderByBuilder() *OrderByBuilder {
	return &OrderByBuilder{}
}

func (b *OrderByBuilder) OrderBy(args ...string) *OrderByBuilder {
	if len(args) == 0 {
		return b
	}
	args[0] = util.BackQuotes(args[0])
	b.args = append(b.args, strings.Join(args, " "))
	return b
}

func (b *OrderByBuilder) OrderByRaw(args ...string) *OrderByBuilder {
	b.args = append(b.args, args...)
	return b
}

func (b *OrderByBuilder) BuildOrderBy() (sqlSegment string) {
	if len(b.args) > 0 {
		sqlSegment = fmt.Sprintf("ORDER BY %s", strings.Join(b.args, ","))
	}
	return
}
