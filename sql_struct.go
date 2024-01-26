package dbgo

//func (db Database) BuildSqlQueryStruct() (sql4prepare string, values []any, err error) {
//	var distinct = db.distinct
//	var fields, bindValuesSelect = db.BuildSelect()
//	tables, err := db.TableBuilder.BuildTable()
//	if err != nil {
//		return sql4prepare, values, err
//	}
//	joins, bindValuesJoin, err := db.BuildJoin()
//	if err != nil {
//		return sql4prepare, values, err
//	}
//	wheres, bindValuesWhere, err := db.BuildWhere()
//	if err != nil {
//		return sql4prepare, values, err
//	}
//	groups, havingS, bindValuesGroup := db.BuildGroup()
//	if err != nil {
//		return sql4prepare, values, err
//	}
//	orderBys := db.BuildOrderBy()
//	pagination, bindValuesPagination := db.BuildPage()
//
//	values = append(values, bindValuesSelect...)
//	values = append(values, bindValuesJoin...)
//	values = append(values, bindValuesWhere...)
//	values = append(values, bindValuesGroup...)
//	values = append(values, bindValuesPagination...)
//
//	if wheres != "" {
//		wheres = fmt.Sprintf("WHERE %s", wheres)
//	}
//	//else {
//	//	// 从struct构建where
//	//}
//	var locks string
//	if db.locking != nil {
//		if *db.locking {
//			locks = "FOR UPDATE"
//		} else {
//			locks = "LOCK IN SHARE MODE"
//		}
//	}
//
//	sql4prepare = util.NamedSprintf("SELECT :distinct :fields FROM :tables :joins :wheres :groups :havings :orderBys :page :locks",
//		distinct, fields, tables, joins, wheres, groups, havingS, orderBys, pagination, locks)
//	return
//}
//func (db Database) BuildSqlExistsStruct() (sql4prepare string, values []any, err error) {
//	sql4prepare, values, err = db.BuildSqlQuery()
//	sql4prepare = fmt.Sprintf("SELECT EXISTS(%s) AS exists", sql4prepare)
//	return
//}
//func (db Database) BuildSqlUpsertStruct(data any, keys []string, columns []string) (sql4prepare string, values []any, err error) {
//	var tmp []string
//	for _, v := range columns {
//		tmp = append(tmp, fmt.Sprintf("`%s`=VALUES(`%s`)", v, v))
//	}
//	return db.buildSqlInsertStruct(data, "", fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", strings.Join(tmp, ", ")))
//}
//func (db Database) BuildSqlInsertUsingStruct(columns []string, b iface.IUnion) (sql4prepare string, values []any, err error) {
//	tables := db.BuildTableOnly4Test()
//	fields := util.Map[string, []string, string](columns, func(s string) string {
//		return fmt.Sprintf("`%s`", s)
//	})
//	prepareQuery, values, err := b.BuildSqlQuery()
//	if err != nil {
//		return sql4prepare, values, err
//	}
//
//	sql4prepare = util.NamedSprintf("INSERT INTO :tables (:fields) (:prepareQuery)", tables, strings.Join(fields, ","), prepareQuery)
//	return sql4prepare, values, err
//}
//func (db Database) BuildSqlInsertOrIgnoreStruct(data any) (sql4prepare string, values []any, err error) {
//	return db.buildSqlInsertStruct(data, "IGNORE")
//}
//func (db Database) BuildSqlInsertStruct(data any) (sql4prepare string, values []any, err error) {
//	return db.buildSqlInsertStruct(data, "")
//}
//func (db Database) buildSqlInsertStruct(data any, ignoreCase string, onDuplicateKeys ...string) (sql4prepare string, values []any, err error) {
//	rfv := reflect.Indirect(reflect.ValueOf(data))
//	var fields []string
//	var valuesPlaceholderArr []string
//	switch rfv.Kind() {
//	case reflect.Struct:
//		keys := rfv.MapKeys()
//		sort.Slice(keys, func(i, j int) bool {
//			return keys[i].String() < keys[j].String()
//		})
//		var valuesPlaceholderTmp []string
//		for _, key := range keys {
//			fields = append(fields, util.BackQuotes(key.String()))
//			valuesPlaceholderTmp = append(valuesPlaceholderTmp, "?")
//			values = append(values, rfv.MapIndex(key).Interface())
//		}
//		valuesPlaceholderArr = append(valuesPlaceholderArr, fmt.Sprintf("(%s)", strings.Join(valuesPlaceholderTmp, ",")))
//	case reflect.Slice:
//		if rfv.Len() == 0 {
//			return
//		}
//		// 先获取到插入字段
//		keys := rfv.Index(0).MapKeys()
//		sort.Slice(keys, func(i, j int) bool {
//			return keys[i].String() < keys[j].String()
//		})
//		for _, key := range keys {
//			fields = append(fields, util.BackQuotes(key.String()))
//		}
//		// 组合插入数据
//		for i := 0; i < rfv.Len(); i++ {
//			var valuesPlaceholderTmp []string
//			for _, key := range keys {
//				valuesPlaceholderTmp = append(valuesPlaceholderTmp, "?")
//				values = append(values, rfv.Index(i).MapIndex(key).Interface())
//			}
//			valuesPlaceholderArr = append(valuesPlaceholderArr, fmt.Sprintf("(%s)", strings.Join(valuesPlaceholderTmp, ",")))
//		}
//	default:
//		err = errors.New("only map(slice) data supported")
//	}
//	tables := db.BuildTableOnly4Test()
//	var onDuplicateKey string
//	if len(onDuplicateKeys) > 0 {
//		onDuplicateKey = onDuplicateKeys[0]
//	}
//	sql4prepare = util.NamedSprintf("INSERT :ignoreCase INTO :tables (:fields) VALUES :placeholder :onDuplicateKey", ignoreCase, tables, strings.Join(fields, ","), strings.Join(valuesPlaceholderArr, " "), onDuplicateKey)
//	return
//}
//func (db Database) BuildSqlUpdateStruct(data any) (sql4prepare string, values []any, err error) {
//	rfv := reflect.Indirect(reflect.ValueOf(data))
//	var updates []string
//	switch rfv.Kind() {
//	case reflect.Map:
//		keys := rfv.MapKeys()
//		sort.Slice(keys, func(i, j int) bool {
//			return keys[i].String() < keys[j].String()
//		})
//		for _, key := range keys {
//			updates = append(updates, fmt.Sprintf("%s = ?", util.BackQuotes(key.String())))
//			values = append(values, rfv.MapIndex(key).Interface())
//		}
//	default:
//		err = errors.New("only map data supported")
//	}
//	tables := db.BuildTableOnly4Test()
//	wheres, binds, err := db.BuildWhere()
//	if err != nil {
//		return sql4prepare, values, err
//	}
//	values = append(values, binds...)
//	if wheres != "" {
//		wheres = fmt.Sprintf("WHERE %s", wheres)
//	}
//	sql4prepare = util.NamedSprintf("UPDATE :tables SET :updates :wheres", tables, strings.Join(updates, ", "), wheres)
//
//	return
//}
//func (db Database) BuildSqlDeleteStruct(id ...int) (sql4prepare string, values []any, err error) {
//	var dbTmp Database
//	if len(id) > 0 {
//		dbTmp = db.Where("id", id[0])
//	} else {
//		dbTmp = db
//	}
//
//	tables := dbTmp.BuildTableOnly4Test()
//	wheres, binds, err := dbTmp.BuildWhere()
//	if err != nil {
//		return sql4prepare, values, err
//	}
//	values = append(values, binds...)
//	if wheres != "" {
//		wheres = fmt.Sprintf("WHERE %s", wheres)
//	}
//	sql4prepare = util.NamedSprintf("DELETE FROM :tables :wheres", tables, wheres)
//	return
//}
//
//// BuildSqlIncrementStruct clause
//// examples
////
////	BuildSqlIncrement("age")
////	BuildSqlIncrement("age", 2)
////	BuildSqlIncrement("age", 3, map[string]any{"name":"John2", "sex": 1})
//func (db Database) BuildSqlIncrementStruct(column string, args ...any) (sql4prepare string, values []any, err error) {
//	return db.buildSqlIncOrDec("+", column, args...)
//}
//
//func (db Database) BuildSqlIncrementEachStruct(data map[string]int, extra ...any) (sql4prepare string, values []any, err error) {
//	return db.buildSqlIncOrDecEach("+", data, extra...)
//}
//
//func (db Database) BuildSqlDecrementStruct(column string, args ...any) (sql4prepare string, values []any, err error) {
//	return db.buildSqlIncOrDec("-", column, args...)
//}
//
//func (db Database) BuildSqlDecrementEachStruct(data map[string]int, extra ...any) (sql4prepare string, values []any, err error) {
//	return db.buildSqlIncOrDecEach("-", data, extra...)
//}
//
//func (db Database) buildSqlIncOrDecStruct(incDec string, column string, args ...any) (sql4prepare string, values []any, err error) {
//	var data = map[string]int{}
//	switch len(args) {
//	case 0:
//		data[column] = 1
//		return db.buildSqlIncOrDecEach(incDec, data)
//	case 1:
//		data[column] = args[0].(int)
//		return db.buildSqlIncOrDecEach(incDec, data)
//	case 2:
//		data[column] = args[0].(int)
//		return db.buildSqlIncOrDecEach(incDec, data, args[1])
//	}
//
//	return
//}
//
//// buildSqlIncOrDecEach specific
//// @incDec +/-
//func (db Database) buildSqlIncOrDecEachStruct(incDec string, data map[string]int, extra ...any) (sql4prepare string, values []any, err error) {
//	var updates []string
//	for k, v := range data {
//		updates = append(updates, fmt.Sprintf("%s = %s %s %v", util.BackQuotes(k), util.BackQuotes(k), incDec, v))
//	}
//
//	//var updates []string
//	if len(extra) > 0 {
//		rfv := reflect.ValueOf(extra[0])
//		keys := rfv.MapKeys()
//		sort.Slice(keys, func(i, j int) bool {
//			return keys[i].String() < keys[j].String()
//		})
//		for _, key := range keys {
//			updates = append(updates, fmt.Sprintf("%s = ?", util.BackQuotes(key.String())))
//			values = append(values, rfv.MapIndex(key).Interface())
//		}
//	}
//
//	tables := db.BuildTableOnly4Test()
//	wheres, binds, err := db.BuildWhere()
//	if err != nil {
//		return sql4prepare, values, err
//	}
//	values = append(values, binds...)
//	if wheres != "" {
//		wheres = fmt.Sprintf("WHERE %s", wheres)
//	}
//	sql4prepare = util.NamedSprintf("UPDATE :tables SET :incDec :updates :wheres", tables, incDec, strings.Join(updates, ", "), wheres)
//
//	return
//}
//
//func (db Database) ToSqlOnlyStruct() string {
//	sql4prepare, _, _ := db.BuildSqlQueryStruct()
//	return sql4prepare
//}
