package main

import (
	"go-webs/dbgo2"
	_ "go-webs/dbgo2/drivers/mysql"
	"log"
)

//	func TrimPrefixAndOr(s string) string {
//		return regexp.MustCompile(`(?i)\s*[and|or]\s*`).ReplaceAllString(s, "")
//	}
type User struct {
	Id   int64  `db:"id,pk"`
	Name string `db:"name"`
}

func main() {
	//a := " and sdffd"
	//log.Println(TrimPrefixAndOr(a))
	//return
	//TestDatabase_ToSql()
	//TestDatabase_ToSql2()
	//TestDatabase_ToSql3()
	//TestDatabase_ToSql4()
	//TestDatabase_ToSqlInsert()
	//TestDatabase_ToSqlUpdate()
	TestDatabase_ToSqlDelete()
}

func db() dbgo2.Database {
	return dbgo2.Open(nil).NewDatabase()
}

func TestDatabase_ToSql() {
	prepare, values, err := db().Table("a").Join("t", "a.id", "t.aid").Select("b").Where("c", 1).OrderBy("id").Limit(10).Page(2).ToSql()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(prepare)
	log.Println(values)
}

func TestDatabase_ToSql2() {
	var user User
	prepare, values, err := db().ToSqlTo(&user)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(prepare)
	log.Println(values)
}

func TestDatabase_ToSql3() {
	var user = User{Id: 1}
	prepare, values, err := db().ToSqlTo(&user)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(prepare)
	log.Println(values)
}

func TestDatabase_ToSql4() {
	var user []User
	prepare, values, err := db().ToSqlTo(&user)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(prepare)
	log.Println(values)
}

func TestDatabase_ToSqlInsert() {
	var user = User{Name: "john"}
	prepare, values, err := db().ToSqlInsert(&user)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(prepare)
	log.Println(values)
}

func TestDatabase_ToSqlUpdate() {
	var user = User{Id: 1, Name: "john"}
	prepare, values, err := db().ToSqlUpdate(&user)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(prepare)
	log.Println(values)
}
func TestDatabase_ToSqlDelete() {
	var user = User{Id: 1, Name: "john"}
	prepare, values, err := db().ToSqlDelete(&user)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(prepare)
	log.Println(values)
}
