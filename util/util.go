package util

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
)

func PtrBool(arg bool) *bool {
	return &arg
}
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
func ToSlice(arg any) []any {
	ref := reflect.Indirect(reflect.ValueOf(arg))
	var res []any
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
func ToSliceAddressable(arg any) []any {
	ref := reflect.Indirect(reflect.ValueOf(arg))
	var res []any
	switch ref.Kind() {
	case reflect.Slice:
		l := ref.Len()
		v := ref.Slice(0, l)
		for i := 0; i < l; i++ {
			res = append(res, v.Index(i).Addr().Interface())
		}
	default:
		res = append(res, ref.Addr().Interface())
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
func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
func GetRandomInt(num int) int {
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
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(fmt.Sprintf(regexp.MustCompile(`:\w+`).ReplaceAllString(format, "%s"), a...), " "))
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

func SortedMapKeys(data any) (cols []string) {
	// 从 map 中获取所有的键，并转换为切片
	keys := reflect.ValueOf(data).MapKeys()

	// 对切片进行排序
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	// 输出排序后的结果
	for _, key := range keys {
		cols = append(cols, key.String())
	}
	return
}
