//下边是一段golang程序,如何在 First() 方法内 获取到调用者 Users 对象

package main

import "log"

type Model struct {
	AAA int
}
type Users struct {
	*Model
	Age int
}

func (u *Model) SetA() {
	u.AAA = 123
}
func (u Users) First() string {
	u.Age = 3
	u.Model.SetA()
	//todo 获取到调用者 Users 对象
	return "First"
}

func main() {
	var u = Users{&Model{}, 3}
	log.Println(u.First())
	log.Println(u.Age)
	log.Println(u.AAA)
}
