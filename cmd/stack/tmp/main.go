package main

import (
	"log"
	"reflect"
)

type rawStruct struct {
	expression string
	binds      []any
}

func main() {
	log.Println(reflect.TypeOf(rawStruct{}).Kind().String())
	log.Println(reflect.TypeOf(rawStruct{}).Name())

}
