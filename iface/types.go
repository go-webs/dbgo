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

type Paginate struct {
	Limit       int                      `json:"limit"`
	Pages       int                      `json:"pages"`
	CurrentPage int                      `json:"currentPage"`
	PrevPage    int                      `json:"prevPage"`
	NextPage    int                      `json:"nextPage"`
	Total       int64                    `json:"total"`
	Data        []map[string]interface{} `json:"data"`
}
