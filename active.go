package dbgo

import (
	"errors"
	"fmt"
	"gitub.com/go-webs/dbgo/iface"
	"gitub.com/go-webs/dbgo/util"
	"reflect"
)

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

	err = db.query(&result, prepare, values...)
	return
}
func (db Database) Find(id int) (result map[string]any, err error) {
	return db.Where("id", id).First()
}
func (db Database) Max(column string) (result []map[string]any, err error) {
	prepare, values, err := db.SelectRaw(fmt.Sprintf("max(%s) as max", util.BackQuotes(column))).BuildSqlQuery()
	if err != nil {
		return result, err
	}

	err = db.query(&result, prepare, values...)
	return
}
func (db Database) Min() (result []map[string]any, err error)   { return }
func (db Database) Avg() (result []map[string]any, err error)   { return }
func (db Database) Count() (result []map[string]any, err error) { return }
func (db Database) Value() (result []map[string]any, err error) { return }
func (db Database) Pluck() (result []map[string]any, err error) { return }
func (db Database) Chunk() (result []map[string]any, err error) { return }
func (db Database) Exists() (result bool, err error)            { return }

func (db Database) Insert(data any, mustFields ...string) (affectedRows int64, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(data))
	switch rfv.Kind() {
	case reflect.Map:
		sql4prepare, values, err := db.BuildSqlInsert(data)
		if err != nil {
			return affectedRows, err
		}
		return db.execute(false, sql4prepare, values...)
	case reflect.Struct:
		err = db.BuildFieldsExecute(data, mustFields...)
		if err != nil {
			return
		}
		return db.Table(data).Insert(db.Datas)
	case reflect.Slice:
		switch rfv.Type().Elem().Kind() {
		case reflect.Map:
			sql4prepare, values, err := db.BuildSqlInsert(data)
			if err != nil {
				return affectedRows, err
			}
			return db.execute(false, sql4prepare, values...)
		case reflect.Struct:
			err = db.BuildFieldsExecute(data, mustFields...)
			if err != nil {
				return
			}
			return db.Table(data).Insert(db.Datas)
		default:
			err = errors.New("data must be map(slice) or struct(slice)")
		}
	default:
		err = errors.New("data must be map(slice) or struct(slice)")
	}
	return
}
func (db Database) InsertOrIgnore(data any) (affectedRows int64, err error) { return }
func (db Database) InsertUsing(columns []string, query iface.IUnion) (affectedRows int64, err error) {
	return
}
func (db Database) Upsert(data any, keys []string, columns []string) (affectedRows int64, err error) {
	return
}

// Update data by where or primary key in data filed
// params
//
//	data 更改的数据,map或者struct类型
//	mustFields 用于struct的0值类型强制更新
//
// 如果没有where条件的情况下,会使用struct的主键(默认id)值作为更新条件,如果都不存在,则禁止更新. 如果是map类型,则必须有where条件
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
		//if !db.WhereBuilder.Exists {
		//	//todo 如果没有where条件,使用主键
		//}
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
		err = errors.New("data must be map(slice) or struct(slice)")
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
		return dbTmp.Update(&values)
	}
	return dbTmp.Insert(&values)
}
func (db Database) Increment() (result []map[string]any, err error)     { return }
func (db Database) Decrement() (result []map[string]any, err error)     { return }
func (db Database) IncrementEach() (result []map[string]any, err error) { return }
func (db Database) DecrementEach() (result []map[string]any, err error) { return }

func (db Database) Delete() (result []map[string]any, err error)   { return }
func (db Database) Truncate() (result []map[string]any, err error) { return }
