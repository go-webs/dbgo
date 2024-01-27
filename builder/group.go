package builder

import (
	"fmt"
	"github.com/go-webs/dbgo/util"
	"strings"
)

type havingStruct struct {
	relation string
	sqlSeg   string
	binds    []any
}
type GroupBuilder struct {
	groups  []string
	havingS []havingStruct
}

func NewGroupBuilder() *GroupBuilder {
	return &GroupBuilder{}
}

// GroupBy clause
// examples
//
//	GroupBy("user_id")
//	GroupBy("user_id", "age")
func (g *GroupBuilder) GroupBy(args ...string) *GroupBuilder {
	if len(args) == 0 {
		return g
	}

	g.groups = append(g.groups,
		util.Map[string, []string, string](args, func(data string) string {
			return util.BackQuotes(data)
		})...)
	return g
}

// GroupByRaw clause
// examples
//
//	GroupByRaw("user_id, age")
func (g *GroupBuilder) GroupByRaw(args ...string) *GroupBuilder {
	g.groups = append(g.groups, args...)
	return g
}

// Having clause
// examples
//
//	Having("count is not null")
//	OrHaving("count",0)
//	Having("count","=",0)
//	HavingRaw("count between ? and ?", [1,2])
//	OrHavingRaw("count between ? and ?", 1, 2)
func (g *GroupBuilder) Having(args ...any) *GroupBuilder      { return g.having("AND", args...) }
func (g *GroupBuilder) OrHaving(args ...any) *GroupBuilder    { return g.having("OR", args...) }
func (g *GroupBuilder) HavingRaw(args ...any) *GroupBuilder   { return g.havingRaw("AND", args...) }
func (g *GroupBuilder) OrHavingRaw(args ...any) *GroupBuilder { return g.havingRaw("OR", args...) }
func (g *GroupBuilder) having(relation string, args ...any) *GroupBuilder {
	if len(args) == 0 {
		return g
	}
	var field = util.BackQuotes(args[0])
	var seg string
	var binds []any
	switch len(args) {
	case 1:
		seg = args[0].(string)
	case 2:
		seg = fmt.Sprintf("%s = ?", field)
		binds = []any{args[1]}
	//case 3 | >3:
	default:
		seg = fmt.Sprintf("%s %s ?", field, args[1])
		binds = []any{args[2]}
	}
	return g.havingRaw(relation, seg, binds)
}
func (g *GroupBuilder) havingRaw(relation string, args ...any) *GroupBuilder {
	switch len(args) {
	// havingRaw("and", "id>age")
	case 1:
		g.havingS = append(g.havingS, havingStruct{relation: relation, sqlSeg: args[0].(string)})
	// havingRaw("and", "id>?", [1])
	case 2:
		g.havingS = append(g.havingS, havingStruct{relation: relation, sqlSeg: args[0].(string), binds: util.ToSlice(args[1])})
	// havingRaw("and", "id between ? and ?", 1, 5)
	default:
		g.havingS = append(g.havingS, havingStruct{relation: relation, sqlSeg: args[0].(string), binds: args[1:]})
	}
	return g
}

func (g *GroupBuilder) BuildGroup() (sqlSegmentGroup, sqlSegmentHaving string, binds []any) {
	if len(g.groups) > 0 {
		sqlSegmentGroup = fmt.Sprintf("GROUP BY %s", strings.Join(
			util.Map[string, []string, string](g.groups, func(data string) string {
				return data
			}), ","))
	}

	var havingArr []string
	for _, v := range g.havingS {
		if v.sqlSeg == "" {
			continue
		}
		if len(havingArr) > 0 {
			havingArr = append(havingArr, fmt.Sprintf("%s %s", v.relation, v.sqlSeg))
		} else {
			havingArr = append(havingArr, v.sqlSeg)
		}
		binds = append(binds, v.binds...)
	}
	if len(havingArr) > 0 {
		sqlSegmentHaving = fmt.Sprintf("HAVING %s", strings.Join(havingArr, " "))
	}
	return
}
