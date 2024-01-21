package main

import (
	"log"
	"reflect"
)

type DB struct{}

type Model struct {
	*DB
	Binder interface{}
}
type Users struct {
	Model `db:"users"`
	Id    int `db:"id,fillable"`
}

func UserModel() *Users {
	return &Users{Model: Model{&DB{}, Users{}}}
}

func (m Model) First() string {
	caller := reflect.ValueOf(m).Field(0).Interface()
	users, ok := caller.(Users)
	if ok {
		log.Println(users)
		// 获取到调用者 Users 对象
		return "First - Got Users"
	}
	return "First"
}

func main() {
	//var a = "a"
	var b = Raw("b")
	print(reflect.TypeOf(b).Name())
	print(reflect.TypeOf(b).PkgPath())

}
