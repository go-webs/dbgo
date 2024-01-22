package util

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
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

func Map[Data any, Datas ~[]Data, Result any](datas Datas, mapper func(Data) Result) []Result {
	results := make([]Result, 0, len(datas))
	for _, data := range datas {
		results = append(results, mapper(data))
	}
	return results
}

func NamedSprintf(format string, a ...any) string {
	str := regexp.MustCompile(`:\w+`).ReplaceAllString(format, "%s")
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(fmt.Sprintf(str, a...), " "))
}

func BackQuotes(arg any) string {
	var tmp []string
	if v, ok := arg.(string); ok {
		split := strings.Split(v, " ")
		split2 := strings.Split(split[0], ".")
		if len(split2) > 1 {
			tmp = append(tmp, fmt.Sprintf("`%s`.`%s`", split2[0], split2[1]))
		} else {
			tmp = append(tmp, fmt.Sprintf("`%s`", split2[len(split2)-1]))
		}
		tmp = append(tmp, split[1:]...)
	}
	return strings.Join(tmp, " ")
}
