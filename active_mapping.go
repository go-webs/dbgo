package dbgo

import (
	"fmt"
	"gitub.com/go-webs/dbgo/util"
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

func (db Database) Insert() (result []map[string]any, err error)         { return }
func (db Database) InsertOrIgnore() (result []map[string]any, err error) { return }
func (db Database) InsertUsing() (result []map[string]any, err error)    { return }
func (db Database) Upsert() (result []map[string]any, err error)         { return }

func (db Database) Update() (result []map[string]any, err error)         { return }
func (db Database) UpdateOrInsert() (result []map[string]any, err error) { return }
func (db Database) Increment() (result []map[string]any, err error)      { return }
func (db Database) Decrement() (result []map[string]any, err error)      { return }
func (db Database) IncrementEach() (result []map[string]any, err error)  { return }
func (db Database) DecrementEach() (result []map[string]any, err error)  { return }

func (db Database) Delete() (result []map[string]any, err error)   { return }
func (db Database) Truncate() (result []map[string]any, err error) { return }
