package dbgo2

import (
	"database/sql"
	"errors"
	"reflect"
)

type Session struct {
	master        *sql.DB
	slave         *sql.DB
	tx            *sql.Tx
	autoSavePoint uint8
}

func (t *Session) Exec(query string, args ...any) (sql.Result, error) {
	if t.tx != nil {
		return t.tx.Exec(query, args...)
	}
	return t.master.Exec(query, args...)
}
func (t *Session) Begin() (err error) {
	if t.tx != nil {
		t.autoSavePoint += 1
		return t.SavePoint(t.autoSavePoint)
	}
	t.tx, err = t.master.Begin()
	return
}
func (t *Session) SavePoint(name any) (err error) {
	_, err = t.tx.Exec("SAVEPOINT ?", name)
	return
}
func (t *Session) RollbackTo(name any) (err error) {
	_, err = t.tx.Exec("ROLLBACK TO ?", name)
	return
}
func (t *Session) Rollback() (err error) {
	if t.autoSavePoint > 0 {
		// decrease in advance whether rollbackTo fail
		currentPoint := t.autoSavePoint
		t.autoSavePoint -= 1
		return t.RollbackTo(currentPoint)
	}
	err = t.tx.Rollback()
	if err != nil {
		return
	}
	t.tx = nil
	return
}
func (t *Session) Commit() (err error) {
	if t.autoSavePoint > 0 {
		t.autoSavePoint -= 1
		return
	}
	err = t.tx.Commit()
	if err != nil {
		return
	}
	t.tx = nil
	return
}
func (t *Session) Trans(closer ...func(*Session) error) (err error) {
	if err = t.Begin(); err != nil {
		return
	}
	for _, v := range closer {
		err = v(t)
		if err != nil {
			return t.Rollback()
		}
	}
	return t.Commit()
}

//	if err = t.Begin(); err != nil {
//		return
//	}
//
//}

func (t *Session) Query(query string, args ...any) (*sql.Rows, error) {
	if t.tx != nil {
		return t.tx.Query(query, args...)
	}
	return t.slave.Query(query, args...)
}
func (t *Session) QueryTo(bind any, query string, args ...any) (err error) {
	var stmt *sql.Stmt
	if t.tx != nil {
		if stmt, err = t.tx.Prepare(query); err != nil {
			return
		}
	} else {
		if stmt, err = t.slave.Prepare(query); err != nil {
			return
		}
	}
	var rows *sql.Rows
	if rows, err = stmt.Query(args...); err != nil {
		return
	}
	return t.rowsToBind(rows, bind)
}
func (t *Session) rowsToBind(rows *sql.Rows, bind any) (err error) {
	rfv := reflect.Indirect(reflect.ValueOf(bind))
	switch rfv.Kind() {
	case reflect.Slice:
		switch rfv.Type().Elem().Kind() {
		case reflect.Map:
			return t.rowsToMap(rows, rfv)
		case reflect.Struct:
			return t.rowsToStruct(rows, rfv)
		default:
			return errors.New("only struct(slice) or map(slice) supported")
		}
	case reflect.Map:
		return t.rowsToMap(rows, rfv)
	case reflect.Struct:
		return t.rowsToStruct(rows, rfv)
	default:
		return errors.New("only struct(slice) or map(slice) supported")
	}
}

//	type binders struct {
//		FieldsTag []string
//		FieldsStruct []string
//	}
type binders struct {
	FieldTag    string
	FieldStruct string
}

func parseBindStruct() (FieldTag []string, FieldStruct []string) {
	return
}
func (t *Session) rowsToStruct(rows *sql.Rows, rfv reflect.Value) error {
	FieldTag, FieldStruct := parseBindStruct()

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// 列的个数
	count := len(columns)

	for rows.Next() {
		// 要先扫描到map, 再做字段比对, 因为这里不确定具体字段数量
		entry, err := t.rowsToMapSingle(rows, columns, count)
		if err != nil {
			return err
		}

		if rfv.Kind() == reflect.Slice {
			rfvItem := reflect.New(rfv.Type().Elem())
			for i, key := range FieldTag {
				if v, ok := entry[key]; ok {
					rfvItem.FieldByName(FieldStruct[i]).Set(reflect.ValueOf(v))
				}
			}
			rfv.Set(reflect.Append(rfv, rfvItem))
		} else {
			for i, key := range FieldTag {
				if v, ok := entry[key]; ok {
					rfv.FieldByName(FieldStruct[i]).Set(reflect.ValueOf(v))
				}
			}
			//rfv.Set(entry)
		}
	}

	return nil
}

func (t *Session) rowsToMapSingle(rows *sql.Rows, columns []string, count int) (entry map[string]any, err error) {
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
		return
	}
	// 一条数据的Map (列名和值的键值对)
	entry = make(map[string]any)

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
	return
}
func (t *Session) rowsToMap(rows *sql.Rows, rfv reflect.Value) error {
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// 列的个数
	count := len(columns)

	for rows.Next() {
		entry, err := t.rowsToMapSingle(rows, columns, count)
		if err != nil {
			return err
		}
		if rfv.Kind() == reflect.Slice {
			rfv.Set(reflect.Append(rfv, reflect.ValueOf(entry)))
		} else {
			rfv.Set(reflect.ValueOf(entry))
		}
	}
	return nil
}
