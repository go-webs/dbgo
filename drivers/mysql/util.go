package mysql

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func BackQuotes(arg any) string {
	var tmp []string
	if v, ok := arg.(string); ok {
		split := strings.Split(v, " ")
		split2 := strings.Split(split[0], ".")
		if len(split2) > 1 {
			if split2[1] == "*" {
				tmp = append(tmp, fmt.Sprintf("`%s`.%s", split2[0], split2[1]))
			} else {
				tmp = append(tmp, fmt.Sprintf("`%s`.`%s`", split2[0], split2[1]))
			}
		} else {
			tmp = append(tmp, fmt.Sprintf("`%s`", split2[len(split2)-1]))
		}
		tmp = append(tmp, split[1:]...)
	}
	return strings.Join(tmp, " ")
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

func NamedSprintf(format string, a ...any) string {
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(fmt.Sprintf(regexp.MustCompile(`:\w+`).ReplaceAllString(format, "%s"), a...), " "))
}
