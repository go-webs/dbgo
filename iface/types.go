package iface

import (
	"reflect"
	"strings"
)

type TypeRaw string

func IsTypeRaw(obj any) bool {
	return reflect.TypeOf(obj).Name() == "TypeRaw"
}

func IsExpression(obj any) (b bool) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	if rfv.Kind() == reflect.String && strings.Contains(rfv.String(), "?") {
		b = true
	}
	return
}
