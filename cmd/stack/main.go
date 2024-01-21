//下边是一段golang程序,如何在 First() 方法内 获取到调用者 Users 对象

package main

import (
	"log"
	"runtime"
)

type Model struct{}
type Users struct {
	Model
}

func (Model) First() string {
	//todo 获取到调用者 Users 对象
	return "First"
}

func main() {
	log.Println(Users{}.First())
}
func resolveCaller() {
	// 根据设定获取错误堆栈信息
	for i := 0; i < 3; i++ {
		// 这里+2是因为,函数调用从封装的函数开始,如果我直接在 main 方法写这个 caller,就是1了
		if funcName, file, line, ok := runtime.Caller(i + 2); ok {
			log.Println(runtime.FuncForPC(funcName).Name(), file, line)
		} else { // 如果堆栈层数没有这么多,就直接中断,不需要再继续获取了
			break
		}
	}
}
