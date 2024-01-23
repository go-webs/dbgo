package dbgo

//type selectStruct struct {
//	selects    []any
//	selectRaws []rawStruct
//}

// Select fields
// args: string | iface.TypeRaw
func (db Database) Select(args ...any) Database {
	db.SelectBuilder.Select(args...)
	return db
}

// AddSelect fields
// args: string | iface.TypeRaw
func (db Database) AddSelect(args ...any) Database {
	db.SelectBuilder.AddSelect(args...)
	return db
}

// SelectRaw fields
// params
//
//	arg: expressions
//	binds: bind values
func (db Database) SelectRaw(arg string, binds ...any) Database {
	//db.selectRaw = append(db.selectRaw, rawStruct{arg, binds})
	db.SelectBuilder.SelectRaw(arg, binds...)
	return db
}

//// BuildSelect fields clause
//func (db Database) BuildSelect() (fields string, binds []any) {
//	return db.SelectBuilder.BuildSelect()
//}

//// BuildSelect fields clause
//func (db selectStruct) buildSelect() (fields string, binds []any) {
//	var tmp []string
//	for _, v := range db.selects {
//		if iface.IsTypeRaw(v) {
//			tmp = append(tmp, string(v.(iface.TypeRaw)))
//		} else {
//			field := v.(string)
//			if strings.ContainsAny(field, " ") {
//				tmp = append(tmp, field)
//			} else {
//				tmp = append(tmp, fmt.Sprintf("`%s`", field))
//			}
//		}
//	}
//	for _, v := range db.selectRaws {
//		tmp = append(tmp, v.expression)
//		binds = append(binds, v.binds...)
//	}
//	if len(tmp) == 0 {
//		return "*", binds
//	}
//	return strings.Join(tmp, ","), binds
//}
