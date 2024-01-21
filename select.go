package dbgo

import (
	"fmt"
	"gitub.com/go-webs/dbgo/iface"
	"strings"
)

type selectStruct struct {
	selects    []any
	selectRaws []rawStruct
}

// Select fields
// args: string | iface.TypeRaw
func (db Database) Select(args ...any) Database {
	db.selects = args
	return db
}

// AddSelect fields
// args: string | iface.TypeRaw
func (db Database) AddSelect(args ...any) Database {
	db.selects = append(db.selects, args...)
	return db
}

// BuildSelect fields clause
func (db Database) BuildSelect() (fields string, binds []any) {
	return db.selectStruct.buildSelect()
}

// BuildSelect fields clause
func (db selectStruct) buildSelect() (fields string, binds []any) {
	var tmp []string
	for _, v := range db.selects {
		if iface.IsTypeRaw(v) {
			tmp = append(tmp, string(v.(iface.TypeRaw)))
		} else {
			field := v.(string)
			if strings.ContainsAny(field, " ") {
				tmp = append(tmp, field)
			} else {
				tmp = append(tmp, fmt.Sprintf("`%s`", field))
			}
		}
	}
	for _, v := range db.selectRaws {
		tmp = append(tmp, v.expression)
		binds = append(binds, v.binds...)
	}
	if len(tmp) == 0 {
		return "*", binds
	}
	return strings.Join(tmp, ","), binds
}
