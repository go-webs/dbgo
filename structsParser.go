package dbgo2

import (
	"database/sql/driver"
	"reflect"
	"slices"
	"strings"
)

func StructsToTableName(rft reflect.Type) (tab string) {
	if field, ok := rft.FieldByName("TableName"); ok {
		if field.Tag.Get("db") != "" {
			tab = field.Tag.Get("db")
		}
	}
	if tab == "" {
		if tn := reflect.New(rft).Elem().MethodByName("TableName"); tn.IsValid() {
			tab = tn.Call(nil)[0].String()
		}
	}
	if tab == "" {
		tab = rft.Name()
	}
	return
}

func StructsParse(obj any) (FieldTag []string, FieldStruct []string, pk string) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Struct:
		return structsParse(rfv, false)
	case reflect.Slice:
		return structsParse(rfv, true)
	default:
		return
	}
}

func structsParse(rfv reflect.Value, isSlice bool) (FieldTag []string, FieldStruct []string, pk string) {
	//rfv := reflect.Indirect(reflect.ValueOf(obj))
	if isSlice {
		rft := rfv.Type().Elem()
		if rft.Kind() == reflect.Struct {
			return structsParse(rfv.Index(0), true)
		}
	} else {
		rft := rfv.Type()
		for i := 0; i < rft.NumField(); i++ {
			field := rft.Field(i)
			if field.Anonymous {
				continue
			}
			tag := field.Tag.Get("db")
			if tag == "-" || tag == "TableName" {
				continue
			}
			if tag == "" {
				//field.Tag = reflect.StructTag("db:" + field.Name)
				FieldStruct = append(FieldStruct, field.Name)
			} else {
				if strings.Contains(tag, ",") {
					tags := strings.Split(tag, ",")
					for _, v := range tags {
						if v == "pk" {
							pk = field.Name
							//pkValue = rfv.FieldByName(field.Name).Addr()
						}
						FieldStruct = append(FieldStruct, tag)
					}
				} else {
					if tag == "pk" {
						pk = field.Name
						//pkValue = rfv.FieldByName(field.Name)
					}
					FieldStruct = append(FieldStruct, tag)
				}
			}
			//else {
			//	FieldStruct = append(FieldStruct, field.Tag.Get("db"))
			//}
			//if field.Tag.Get("pk") == "true" {
			//	pk = field.Name
			//	pkValue = rfv.FieldByName(field.Name)
			//}
			FieldTag = append(FieldTag, field.Name)
		}
	}
	return
}

func StructToSelects(obj any) []string {
	tag, fieldStruct, _ := StructsParse(obj)
	if len(tag) > 0 {
		return tag
	} else {
		return fieldStruct
	}
}

func structDataToMap(rfv reflect.Value, tag, fieldStruct []string, mustFields ...string) (data map[string]any, err error) {
	data = make(map[string]any)
	for i, fieldName := range fieldStruct {
		field := rfv.FieldByName(fieldName)
		if (field.Kind() == reflect.Ptr && field.IsNil()) || (field.IsZero() && !slices.Contains(mustFields, tag[i])) {
			continue
		}
		var rfvVal = field.Interface()
		if v, ok := rfvVal.(driver.Valuer); ok {
			var value driver.Value
			value, err = v.Value()
			if err != nil {
				return
			}
			data[tag[i]] = value
		} else {
			data[tag[i]] = rfvVal
		}
	}
	return
}

func StructToDelete(obj any) (pk string, pkValue any) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	if rfv.Kind() == reflect.Struct {
		_, _, pk = StructsParse(obj)
		if pk != "" {
			pkValue = rfv.FieldByName(pk).Interface()
		}
	}
	return
}

func StructsToInsert(obj any, mustFields ...string) (datas []map[string]any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Struct:
		tag, fieldStruct, _ := StructsParse(obj)
		var data map[string]any
		data, err = structDataToMap(rfv, tag, fieldStruct, mustFields...)
		if err != nil {
			return
		}
		datas = append(datas, data)
	case reflect.Slice:
		tag, fieldStruct, _ := StructsParse(obj)
		for i := 0; i < rfv.Len(); i++ {
			var data map[string]any
			data, err = structDataToMap(rfv.Index(i), tag, fieldStruct, mustFields...)
			if err != nil {
				return
			}
			datas = append(datas, data)
		}
	default:
		return
	}
	return
}

func StructToUpdate(obj any, mustFields ...string) (data map[string]any, pk string, pkValue any, err error) {
	tag, fieldStruct, pkCol := StructsParse(obj)
	if len(tag) > 0 {
		data = make(map[string]any)
		rfv := reflect.Indirect(reflect.ValueOf(obj))
		data, err = structDataToMap(rfv, tag, fieldStruct, mustFields...)
		if err != nil {
			return
		}
		if pkCol != "" {
			pk = pkCol
			pkValue = rfv.FieldByName(pk).Interface()
		}
	}

	return
}
