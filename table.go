package dbgo

//type tableStruct struct {
//	table  any
//	as     string
//	prefix string
//}

// Table select from table as a
// param
//
//	arg: table string | struct bindings
//	as: table name alias
func (db Database) Table(table any, as ...string) Database {
	db.TableBuilder.Table(table, as...)
	return db
}
func (db Database) BuildTable() string {
	segment, _ := db.TableBuilder.BuildTable()
	return segment
}

//// BuildTable for table name build
//func (db Database) BuildTable(ts ...tableStruct) (tab string) {
//	table := db.tableStruct
//	if len(ts) > 0 {
//		table = ts[0]
//	}
//	return table.buildTable()
//}
//
//func newTable(prefix string, table any, as ...string) tableStruct {
//	var ts = tableStruct{
//		table:  table,
//		prefix: prefix,
//	}
//	if len(as) > 0 {
//		ts.as = as[0]
//	}
//	return ts
//}
//
//// buildTable name
//func (ts tableStruct) buildTable() (tab string) {
//	rfv := reflect.Indirect(reflect.ValueOf(ts.table))
//	switch rfv.Kind() {
//	case reflect.String:
//		tab = fmt.Sprintf("`%s%s`", ts.prefix, ts.table)
//	case reflect.Struct:
//		tab = ts.buildTableName(rfv.Type())
//	case reflect.Slice:
//		if rfv.Type().Elem().Kind() == reflect.Struct {
//			tab = ts.buildTableName(rfv.Type().Elem())
//		} else {
//			ts.Err = errors.New("table param must be string or struct(slice) bind with 1 or 2 params")
//		}
//	default:
//		ts.Err = errors.New("table must string | struct | slice")
//	}
//	return strings.TrimSpace(fmt.Sprintf("%s %s", tab, ts.as))
//}
//
//func (ts tableStruct) buildTableName(rft reflect.Type) (tab string) {
//	if field, ok := rft.FieldByName("TableName"); ok {
//		if field.Tag.Get("db") != "" {
//			tab = field.Tag.Get("db")
//		}
//	}
//	if tab == "" {
//		tab = fmt.Sprintf("`%s%s`", ts.prefix, rft.Name())
//	} else {
//		tab = fmt.Sprintf("`%s`", tab)
//	}
//	return
//}
