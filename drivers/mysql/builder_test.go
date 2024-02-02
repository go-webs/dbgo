package mysql

import (
	"go-webs/dbgo2"
	"log"
	"testing"
)

type User struct {
	Id   int64 `id,pk`
	Name string
}

func db() dbgo2.Database {
	return dbgo2.Open(nil).NewDatabase()
}

func TestDatabase_ToSql(t *testing.T) {
	prepare, values, err := db().Table("a").Select("b").Where("c", 1).OrderBy("id").Limit(10).Page(2).ToSql()
	if err != nil {
		log.Fatal(err.Error())
	}
	t.Log(prepare)
	t.Log(values)
}
func TestDatabase_ToSqlInsert(t *testing.T) {
	prepare, values, err := db().ToSqlInsert()
	if err != nil {
		log.Fatal(err.Error())
	}
	t.Log(prepare)
	t.Log(values)
}
