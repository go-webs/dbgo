package dbgo

import (
	"math/rand"
	"reflect"
	"time"
)

func PtrString(arg string) *string {
	return &arg
}
func PtrInt(arg int) *int {
	return &arg
}
func PtrInt8(arg int8) *int8 {
	return &arg
}
func PtrInt16(arg int16) *int16 {
	return &arg
}
func PtrInt64(arg int64) *int64 {
	return &arg
}
func PtrFloat64(arg float64) *float64 {
	return &arg
}
func PtrTime(arg time.Time) *time.Time {
	return &arg
}
func ToSlice(arg interface{}) []interface{} {
	ref := reflect.Indirect(reflect.ValueOf(arg))
	var res []interface{}
	switch ref.Kind() {
	case reflect.Slice:
		l := ref.Len()
		v := ref.Slice(0, l)
		for i := 0; i < l; i++ {
			res = append(res, v.Index(i).Interface())
		}
	default:
		res = append(res, ref.Interface())
	}
	return res
}
func SliceContains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func getRandomInt(num int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(num)
}
