package dbgo

import "errors"

func (jc *TypeJoinOnCondition) onClause(relation, column1 string, column2OrOperator ...string) IJoin {
	if len(column2OrOperator) == 0 {
		return jc
	}
	var operator = "="
	var column2 = column2OrOperator[0]
	if len(column2OrOperator) == 2 {
		operator = column2OrOperator[0]
		column2 = column2OrOperator[1]
	}
	jc.Conditions = append(jc.Conditions, TypeJoinOnConditionItem{
		Relation: relation,
		Column1:  column1,
		Operator: operator,
		Column2:  column2,
	})
	return jc
}
func (jc *TypeJoinOnCondition) On(column1 string, operatorOrColumn2 ...string) IJoin {
	return jc.onClause("AND", column1, operatorOrColumn2...)
}
func (jc *TypeJoinOnCondition) OrOn(column1 string, operatorOrColumn2 ...string) IJoin {
	return jc.onClause("OR", column1, operatorOrColumn2...)
}

func (jc *JoinClause) join(joinType string, table any, argOrFn ...any) *JoinClause {
	var tab TableClause
	switch table.(type) {
	case string:
		tab.Tables = table
	case TableClause:
		tab = table.(TableClause)
	case IBuilder:
		jc.JoinItems = append(jc.JoinItems, TypeJoinSub{table.(IBuilder)})
		return jc
	}

	switch len(argOrFn) {
	case 1:
		if v, ok := argOrFn[0].(func(IJoin)); ok {
			jc.JoinItems = append(jc.JoinItems, TypeJoinOn{
				TableClause: tab,
				OnClause:    v,
				Type:        joinType,
			})
		}
	case 2:
		jc.JoinItems = append(jc.JoinItems, TypeJoinStandard{
			TableClause: tab,
			Column1:     argOrFn[0].(string),
			Operator:    "=",
			Column2:     argOrFn[1].(string),
			Type:        joinType,
		})
	case 3:
		jc.JoinItems = append(jc.JoinItems, TypeJoinStandard{
			TableClause: tab,
			Column1:     argOrFn[0].(string),
			Operator:    argOrFn[1].(string),
			Column2:     argOrFn[2].(string),
			Type:        joinType,
		})
	default:
		jc.Err = errors.New("join args error")
	}
	return jc
}
