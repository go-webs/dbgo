package builder

import "reflect"

type BuildValue struct {
	Key   string
	Value any
}

type IBuilder interface {
	//Get(bindStructSeg any)
	Build() (sqlSegment string, bindValues []any, err error)
}

type rawStruct struct {
	expression string
	binds      []any
}

func isRawStruct(arg any) bool {
	return reflect.TypeOf(arg).Name() == "rawStruct"
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
