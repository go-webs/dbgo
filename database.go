package dbgo

import (
	"database/sql"
	"errors"
	"gitub.com/go-webs/dbgo/builder"
	"gitub.com/go-webs/dbgo/util"
	"reflect"
)

type transaction struct {
	tx *sql.Tx
}

func (db *transaction) Begin(sd *sql.DB) (err error) {
	db.tx, err = sd.Begin()
	return
}

// Rollback ...
func (db *transaction) Rollback() (err error) {
	err = db.tx.Rollback()
	db.tx = nil
	return
}

// Commit ...
func (db *transaction) Commit() (err error) {
	err = db.tx.Commit()
	db.tx = nil
	return
}

type Database struct {
	*DbGo
	*transaction
	distinct string
	locking  *bool // Pessimistic Locking

	*builder.TableBuilder
	*builder.SelectBuilder
	*builder.JoinBuilder
	*builder.WhereBuilder
	*builder.GroupBuilder
	*builder.OrderByBuilder
	*builder.PageBuilder
}

func newDatabase(dg *DbGo) *Database {
	return &Database{
		DbGo:           dg,
		TableBuilder:   builder.NewTableBuilder(dg.Cluster.Prefix),
		SelectBuilder:  builder.NewSelectBuilder(),
		JoinBuilder:    builder.NewJoinBuilder(dg.Cluster.Prefix),
		WhereBuilder:   builder.NewWhereBuilder(),
		GroupBuilder:   builder.NewGroupBuilder(),
		OrderByBuilder: builder.NewOrderByBuilder(),
		PageBuilder:    builder.NewPageBuilder(),
	}
}

// Distinct for distinct
func (db Database) Distinct() Database {
	db.distinct = "DISTINCT"
	return db
}

// SharedLock 4 select ... locking in share mode
func (db Database) SharedLock() Database {
	db.locking = util.PtrBool(false)
	return db
}

// LockForUpdate 4 select ... for update
func (db Database) LockForUpdate() Database {
	db.locking = util.PtrBool(true)
	return db
}

// Begin ...
func (db Database) Begin() (err error) {
	err = db.transaction.Begin(db.MasterDB())
	return
}

// Transaction ...
func (db Database) Transaction(closers ...func(db Database) error) (err error) {
	err = db.Begin()
	if err != nil {
		return err
	}
	for _, closer := range closers {
		err = closer(db)
		if err != nil {
			return db.Rollback()
		}
	}
	return db.Commit()
}

func (db Database) query(obj any, sql4prepare string, binds ...any) (err error) {
	var prepare *sql.Stmt
	if db.tx == nil {
		prepare, err = db.SlaveDB().Prepare(sql4prepare)
	} else {
		prepare, err = db.tx.Prepare(sql4prepare)
	}
	defer prepare.Close()
	db.recordSqlLog(sql4prepare, binds...)

	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Slice:
		switch reflect.Indirect(reflect.New(rfv.Type().Elem())).Kind() {
		//case reflect.Struct:
		//	return db.scanStruct(rfv, prepare, binds...)
		case reflect.Map:
			return db.scanMap(rfv, prepare, binds...)
		default:
			return errors.New("unsorted obj")
		}
	//case reflect.Struct:
	//	return db.scanStruct(rfv, prepare, args...)
	case reflect.Map:
		return db.scanMap(rfv, prepare, binds...)
	default:
		return errors.New("unsorted obj")
	}
}
func (db Database) scanMap(rfv reflect.Value, prepare *sql.Stmt, args ...any) error {
	rows, err := prepare.Query(args...)
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// 列的个数
	count := len(columns)

	for rows.Next() {
		// 一条数据的各列的值（需要指定长度为列的个数，以便获取地址）
		values := make([]any, count)
		// 一条数据的各列的值的地址
		valPointers := make([]any, count)
		// 获取各列的值的地址
		for i := 0; i < count; i++ {
			valPointers[i] = &values[i]
		}
		// 获取各列的值，放到对应的地址中
		err = rows.Scan(valPointers...)
		if err != nil {
			return err
		}
		// 一条数据的Map (列名和值的键值对)
		entry := make(map[string]any)

		// Map 赋值
		for i, col := range columns {
			var v any
			// 值复制给val(所以Scan时指定的地址可重复使用)
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				// 字符切片转为字符串
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		if rfv.Kind() == reflect.Slice {
			rfv.Set(reflect.Append(rfv, reflect.ValueOf(entry)))
		} else {
			rfv.Set(reflect.ValueOf(entry))
		}
	}
	return nil
}

//func (db Database) scanStruct(rfv reflect.Value, prepare *sql.Stmt, args ...any) error {
//	dbFields, structFields, structRft, err := db.getFieldsQuery(rfv)
//	if err != nil {
//		return err
//	}
//
//	rows, err := prepare.Query(args...)
//	if err != nil {
//		return err
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		fields := make([]any, len(dbFields))
//		entry := reflect.New(structRft).Elem()
//		for i, v := range structFields {
//			field := entry.FieldByName(v)
//			fields[i] = field.Addr().Interface()
//		}
//		if err := rows.Scan(fields...); err != nil {
//			return err
//		}
//
//		if rfv.Kind() == reflect.Slice {
//			rfv.Set(reflect.Append(rfv, entry))
//		} else {
//			rfv.Set(entry)
//		}
//	}
//
//	return nil
//}
