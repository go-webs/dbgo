package dbgo

import (
	"errors"
	"fmt"
	"github.com/go-webs/dbgo/iface"
	"github.com/go-webs/dbgo/util"
	"math"
	"reflect"
	"slices"
)

func (db Database) To(binds any, mustFields4where ...string) (err error) {
	rfv := reflect.Indirect(reflect.ValueOf(binds))
	var dbTmp = db.Table(binds)
	switch rfv.Kind() {
	case reflect.Struct:
		rft := rfv.Type()
		err = dbTmp.BuildFieldsQuery(rft)
		if err != nil {
			return
		}
		for i := 0; i < len(dbTmp.Bindery.FieldsStruct); i++ {
			dbTmp = dbTmp.Select(dbTmp.Bindery.FieldsTag[i])
			if len(mustFields4where) > 0 && slices.Contains(mustFields4where, dbTmp.Bindery.FieldsTag[i]) ||
				rfv.FieldByName(dbTmp.Bindery.FieldsStruct[i]).Kind() == reflect.Ptr && !rfv.FieldByName(dbTmp.Bindery.FieldsStruct[i]).IsNil() ||
				!rfv.FieldByName(dbTmp.Bindery.FieldsStruct[i]).IsZero() {
				dbTmp = dbTmp.Where(dbTmp.Bindery.FieldsTag[i], rfv.FieldByName(dbTmp.Bindery.FieldsStruct[i]).Interface())
			}
		}
		dbTmp = dbTmp.Limit(1)
	case reflect.Slice:
		if rfv.Type().Elem().Kind() != reflect.Struct {
			err = errors.New("binds must be struct(slice)")
			return
		}
		rft := rfv.Type()
		err = dbTmp.BuildFieldsQuery(rft.Elem())
		if err != nil {
			return
		}
		for _, v := range dbTmp.Bindery.FieldsTag {
			dbTmp = dbTmp.Select(v)
		}
	default:
		err = errors.New("binds must be struct(slice)")
		return
	}
	sql4prepare, values, err2 := dbTmp.BuildSqlQuery()
	if err2 != nil {
		return err2
	}
	return dbTmp.query(binds, sql4prepare, values...)
}
func (db Database) Get() (result []map[string]any, err error) {
	prepare, values, err := db.BuildSqlQuery()
	if err != nil {
		return result, err
	}

	err = db.query(&result, prepare, values...)
	return
}
func (db Database) First() (result map[string]any, err error) {
	prepare, values, err := db.Limit(1).BuildSqlQuery()
	if err != nil {
		return result, err
	}

	result = map[string]any{}
	err = db.query(&result, prepare, values...)
	return
}
func (db Database) Find(id int64) (result map[string]any, err error) {
	return db.Where("id", id).First()
}
func (db Database) Max(column string) (result float64, err error) {
	err = db.aggregation(&result, util.BackQuotes(column), "max")
	return
}
func (db Database) Min(column string) (result float64, err error) {
	err = db.aggregation(&result, util.BackQuotes(column), "min")
	return
}
func (db Database) Avg(column string) (result float64, err error) {
	err = db.aggregation(&result, util.BackQuotes(column), "avg")
	return
}
func (db Database) Count() (result int64, err error) {
	err = db.aggregation(&result, "*", "count")
	return
}
func (db Database) Value(column string) (result any, err error) {
	first, err := db.Select(column).Limit(1).First()
	if err != nil {
		return result, err
	}
	if v, ok := first[column]; ok {
		result = v
	}
	return
}
func (db Database) Pluck(field string, fieldKey ...string) (result any, err error) {
	if len(fieldKey) > 0 {
		get, err := db.Select(field, fieldKey[0]).Get()
		if err != nil {
			return result, err
		}
		var tmp = make(map[any]any)
		for _, v := range get {
			tmp[v[fieldKey[0]]] = v[field]
		}
		result = tmp
	} else {
		get, err := db.Select(field).Get()
		if err != nil {
			return result, err
		}
		var tmp []any
		for _, v := range get {
			tmp = append(tmp, v[field])
		}
		result = tmp
	}
	return
}
func (db Database) DoesntExist() (result bool, err error) {
	result, err = db.Exists()
	return !result, err
}
func (db Database) Exists() (result bool, err error) {
	prepare, values, err := db.BuildSqlExists()
	if err != nil {
		return result, err
	}
	err = db.queryRow(&result, prepare, values...)
	return
}
func (db Database) Chunk(limit int, callback func(dataList []map[string]any) error) error {
	count, err2 := db.Count()
	if err2 != nil {
		return err2
	}
	pages := int(math.Ceil(float64(count) / float64(limit)))
	for i := 1; i <= pages; i++ {
		result, err := db.Limit(limit).Page(i).Get()
		if err != nil {
			return err
		}
		err = callback(result)
		if err != nil {
			return err
		}
	}
	return nil
}

type Lazy struct {
	Database
}

func (db Database) Lazy() Lazy {
	return Lazy{db}
}
func (db Lazy) Each(callback func(data map[string]any) error) error {
	count, err2 := db.Count()
	if err2 != nil {
		return err2
	}
	for i := 1; i <= int(count); i++ {
		result, err := db.Limit(100).Page(i).Get()
		if err != nil {
			return err
		}
		for j := 0; j < len(result); j++ {
			err = callback(result[j])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (db Database) Insert(data any, mustFields ...string) (affectedRows int64, err error) {
	return db.insert(false, false, data, mustFields...)
}
func (db Database) InsertGetId(data any, mustFields ...string) (lastInsertId int64, err error) {
	return db.insert(true, false, data, mustFields...)
}
func (db Database) InsertOrIgnore(data any, mustFields ...string) (affectedRows int64, err error) {
	return db.insert(false, true, data, mustFields...)
}
func (db Database) InsertUsing(columns []string, query iface.IUnion) (affectedRows int64, err error) {
	prepare, values, err := db.BuildSqlInsertUsing(columns, query)
	if err != nil {
		return affectedRows, err
	}
	return db.execute(false, prepare, values...)
}
func (db Database) Upsert(data any, keys []string, columns []string) (affectedRows int64, err error) {
	prepare, values, err := db.BuildSqlUpsert(data, keys, columns)
	if err != nil {
		return affectedRows, err
	}
	return db.execute(false, prepare, values...)
}

// Update data by where or primary key in data filed
// params
//
//	data 更改的数据,map或者struct类型
//	mustFields 用于struct的0值类型强制更新
//
// 优先使用struct的主键(默认id)值作为更新条件,如果不存在且没有where条件,则禁止更新
// examples
//
//	db.Where("id",1).Update({"name":"david"})	// update users set name=? where id=1
//	db.Update(Users{Id:1, Name:"David"})	// update users set name=? where id=1
//	db.Update(Users{Id:1, Name:"David", Votes:0}, "votes")	// update users set name=?, votes=? where id=?
func (db Database) Update(data any, mustFields ...string) (affectedRows int64, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(data))
	switch rfv.Kind() {
	case reflect.Map:
		sql4prepare, values, err := db.BuildSqlUpdate(data)
		if err != nil {
			return affectedRows, err
		}
		return db.execute(false, sql4prepare, values...)
	case reflect.Struct:
		err = db.BuildFieldsExecute(data, mustFields...)
		if err != nil {
			return
		}
		//dbTmp := db
		//if !db.WhereBuilderNew.Exists {
		//	//todo 如果没有where条件,使用主键
		//	dbTmp = db.Where(db.Bindery.PrimaryKey, db.Bindery.PrimaryKeyValue)
		//}
		if db.Bindery.PrimaryKey != "" && db.Bindery.PrimaryKeyValue != nil {
			return db.Table(data).Where(db.Bindery.PrimaryKey, db.Bindery.PrimaryKeyValue).Update(db.Datas[0])
		}
		return db.Table(data).Update(db.Datas)
	//case reflect.Slice:
	//	switch rfv.Type().Elem().Kind() {
	//	case reflect.Map:
	//		sql4prepare, values, err := db.BuildSqlUpdate(data)
	//		if err != nil {
	//			return affectedRows, err
	//		}
	//		return db.execute(false, sql4prepare, values...)
	//	case reflect.Struct:
	//		err = db.BuildFieldsExecute(data, mustFields...)
	//		if err != nil {
	//			return
	//		}
	//		return db.Table(data).Update(db.Datas)
	//	default:
	//		err = errors.New("data must be map(slice) or struct(slice)")
	//	}
	default:
		err = errors.New("data must be map or struct")
	}
	return
}
func (db Database) UpdateOrInsert(attributes, values map[string]any) (affectedRows int64, err error) {
	dbTmp := db.Where(attributes)
	var exists bool
	if exists, err = dbTmp.Exists(); err != nil {
		return
	}
	if exists {
		return dbTmp.Update(values)
	}
	return dbTmp.Insert(values)
}
func (db Database) Increment(column string, args ...any) (affectedRows int64, err error) {
	prepare, values, err := db.BuildSqlIncrement(column, args...)
	if err != nil {
		return affectedRows, err
	}
	return db.execute(false, prepare, values...)
}
func (db Database) Decrement(column string, args ...any) (affectedRows int64, err error) {
	prepare, values, err := db.BuildSqlDecrement(column, args...)
	if err != nil {
		return affectedRows, err
	}
	return db.execute(false, prepare, values...)
}
func (db Database) IncrementEach(data map[string]any, extra ...any) (affectedRows int64, err error) {
	prepare, values, err := db.BuildSqlIncrementEach(data, extra...)
	if err != nil {
		return affectedRows, err
	}
	return db.execute(false, prepare, values...)
}
func (db Database) DecrementEach(data map[string]any, extra ...any) (affectedRows int64, err error) {
	prepare, values, err := db.BuildSqlDecrementEach(data, extra...)
	if err != nil {
		return affectedRows, err
	}
	return db.execute(false, prepare, values...)
}

func (db Database) Delete(id ...int64) (affectedRows int64, err error) {
	prepare, values, err := db.BuildSqlDelete(id...)
	if err != nil {
		return affectedRows, err
	}
	return db.execute(false, prepare, values...)
}
func (db Database) Truncate(obj ...any) (affectedRows int64, err error) {
	var table string
	var dbTmp = db
	if len(obj) > 0 {
		dbTmp = db.Table(obj[0])
	}
	table, err = dbTmp.BuildTable()
	if err != nil {
		return
	}
	return db.execute(false, fmt.Sprintf("TRUNCATE TABLE %s", table))
}

func (db Database) Paginate(obj ...any) (result iface.Paginate, err error) {
	var count int64
	dbTmp := db
	if len(obj) > 0 {
		dbTmp = db.Table(obj[0])
	}
	count, err = dbTmp.Count()
	if err != nil || count == 0 {
		return
	}
	limit, _, page := db.PageBuilder.GetPagination()
	if limit == 0 {
		limit = 15
	}
	if page == 0 {
		page = 1
	}

	var res []map[string]any
	res, err = dbTmp.Limit(limit).Page(page).Get()
	if err != nil {
		return
	}

	result.Total = count
	result.Data = res
	result.Limit = limit
	result.Pages = int(math.Ceil(float64(count) / float64(limit)))
	result.CurrentPage = page
	result.PrevPage = page - 1
	result.NextPage = page + 1
	if page == 1 {
		result.PrevPage = 1
	}
	if page == result.Pages {
		result.NextPage = result.Pages
	}
	return
}

func (db Database) aggregation(bind any, column string, agg string) (err error) {
	//prepare, values, err2 := db.SelectRaw(fmt.Sprintf("%s(%s) as %s", agg, column, agg)).BuildSqlQuery()
	//if err2 != nil {
	//	return err2
	//}
	segment, err2 := db.BuildTable()
	if err2 != nil {
		return err2
	}
	prepare := fmt.Sprintf("SELECT %s(%s) as %s FROM %s", agg, column, agg, segment)
	return db.queryRow(bind, prepare)
}
func (db Database) insert(returnLastInsertId, ignore bool, data any, mustFields ...string) (affectedRows int64, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(data))
	switch rfv.Kind() {
	case reflect.Map:
		if rfv.Len() == 0 {
			return
		}
		var sql4prepare string
		var values []any
		if ignore {
			sql4prepare, values, err = db.BuildSqlInsertOrIgnore(data)
		} else {
			sql4prepare, values, err = db.BuildSqlInsert(data)
		}
		if err != nil {
			return affectedRows, err
		}
		return db.execute(returnLastInsertId, sql4prepare, values...)
	case reflect.Struct:
		err = db.BuildFieldsExecute(data, mustFields...)
		if err != nil {
			return
		}
		return db.Table(data).insert(returnLastInsertId, ignore, db.Datas, mustFields...)
	case reflect.Slice:
		if rfv.Len() == 0 {
			return
		}
		switch rfv.Type().Elem().Kind() {
		case reflect.Map:
			var sql4prepare string
			var values []any
			if ignore {
				sql4prepare, values, err = db.BuildSqlInsertOrIgnore(data)
			} else {
				sql4prepare, values, err = db.BuildSqlInsert(data)
			}
			if err != nil {
				return affectedRows, err
			}
			return db.execute(returnLastInsertId, sql4prepare, values...)
		case reflect.Struct:
			err = db.BuildFieldsExecute(data, mustFields...)
			if err != nil {
				return
			}
			return db.Table(data).insert(returnLastInsertId, ignore, db.Datas, mustFields...)
		default:
			err = errors.New("data must be map(slice) or struct(slice)")
		}
	default:
		err = errors.New("data must be map(slice) or struct(slice)")
	}
	return
}
