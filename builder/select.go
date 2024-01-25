package builder

import (
	"gitub.com/go-webs/dbgo/iface"
	"gitub.com/go-webs/dbgo/util"
	"strings"
)

type SelectBuilder struct {
	selects []any
}

func NewSelectBuilder() *SelectBuilder {
	return &SelectBuilder{}
}

func (ss *SelectBuilder) Select(selects ...any) {
	ss.selects = append(ss.selects, selects...)
}

func (ss *SelectBuilder) AddSelect(selects ...any) {
	ss.selects = append(ss.selects, selects...)
}

func (ss *SelectBuilder) SelectRaw(selects string, binds ...any) {
	ss.AddSelect(rawStruct{selects, binds})
}

//func (ss *SelectBuilder) BuildTableOnly4Test() {
//
//}

// BuildSelect fields clause
func (ss *SelectBuilder) BuildSelect() (fields string, binds []any) {
	var tmp []string
	for _, v := range ss.selects {
		if iface.IsTypeRaw(v) {
			tmp = append(tmp, string(v.(iface.TypeRaw)))
		} else if isRawStruct(v) {
			tmp = append(tmp, v.(rawStruct).expression)
			binds = append(binds, v.(rawStruct).binds...)
		} else {
			field := v.(string)
			if strings.ContainsAny(field, " ") {
				tmp = append(tmp, field)
			} else {
				tmp = append(tmp, util.BackQuotes(field))
			}
		}
	}

	if len(tmp) == 0 {
		return "*", binds
	}
	return strings.Join(tmp, ","), binds
}
