package main

import (
	"go-webs/dbgo2"
	_ "go-webs/dbgo2/drivers/mysql"
	"log"
)

func main() {
	TestDatabase_ToSql()
}

func db() dbgo2.Database {
	return dbgo2.Open(nil).NewDatabase()
}
func TestDatabase_ToSql() {
	prepare, values, err := db().Table("a").Select("b").Where("c", 1).OrderBy("id").Limit(10).Page(2).ToSql()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(prepare)
	log.Println(values)
}
