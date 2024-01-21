package iface

import "reflect"

type TypeRaw string

func IsTypeRaw(obj interface{}) bool {
	return reflect.TypeOf(obj).Name() == "TypeRaw"
}
