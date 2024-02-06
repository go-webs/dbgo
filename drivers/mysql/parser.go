package mysql

import (
	"errors"
	"fmt"
	"go-webs/dbgo2"
	"reflect"
	"sort"
	"strings"
)

func (b Builder) buildTableName(rft reflect.Type, prefix string) (tab string) {
	return BackQuotes(fmt.Sprintf("%s%s", prefix, dbgo2.StructsToTableName(rft)))
}

func (b Builder) toSqlInsert(c *dbgo2.Context, data any, ignoreCase string, onDuplicateKeys ...string) (sql4prepare string, values []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(data))
	var fields []string
	var valuesPlaceholderArr []string
	switch rfv.Kind() {
	case reflect.Map:
		keys := rfv.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		var valuesPlaceholderTmp []string
		for _, key := range keys {
			fields = append(fields, BackQuotes(key.String()))
			valuesPlaceholderTmp = append(valuesPlaceholderTmp, "?")
			values = append(values, rfv.MapIndex(key).Interface())
		}
		valuesPlaceholderArr = append(valuesPlaceholderArr, fmt.Sprintf("(%s)", strings.Join(valuesPlaceholderTmp, ",")))
	case reflect.Slice:
		if rfv.Len() == 0 {
			return
		}
		if rfv.Type().Elem().Kind() == reflect.Map {
			// 先获取到插入字段
			keys := rfv.Index(0).MapKeys()
			sort.Slice(keys, func(i, j int) bool {
				return keys[i].String() < keys[j].String()
			})
			for _, key := range keys {
				fields = append(fields, BackQuotes(key.String()))
			}
			// 组合插入数据
			for i := 0; i < rfv.Len(); i++ {
				var valuesPlaceholderTmp []string
				for _, key := range keys {
					valuesPlaceholderTmp = append(valuesPlaceholderTmp, "?")
					values = append(values, rfv.Index(i).MapIndex(key).Interface())
				}
				valuesPlaceholderArr = append(valuesPlaceholderArr, fmt.Sprintf("(%s)", strings.Join(valuesPlaceholderTmp, ",")))
			}
		} else {
			err = errors.New("only map(slice) data supported")
			return
		}
	default:
		err = errors.New("only map(slice) data supported")
		return
	}
	if err != nil {
		return
	}
	var onDuplicateKey string
	if len(onDuplicateKeys) > 0 {
		onDuplicateKey = onDuplicateKeys[0]
	}
	var tables string
	tables, _, err = b.ToSqlTable(c)
	if err != nil {
		return
	}
	sql4prepare = NamedSprintf("INSERT :ignoreCase INTO :tables (:fields) VALUES :placeholder :onDuplicateKey", ignoreCase, tables, strings.Join(fields, ","), strings.Join(valuesPlaceholderArr, ","), onDuplicateKey)
	return
}

func (b Builder) toSqlUpdate(c *dbgo2.Context, data any) (sql4prepare string, values []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(data))
	var updates []string
	switch rfv.Kind() {
	case reflect.Map:
		keys := rfv.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		for _, key := range keys {
			updates = append(updates, fmt.Sprintf("%s = ?", BackQuotes(key.String())))
			values = append(values, rfv.MapIndex(key).Interface())
		}
	default:
		err = errors.New("only map data supported")
		return
	}
	var tables string
	tables, _, err = b.ToSqlTable(c)
	if err != nil {
		return
	}
	wheres, binds, err := b.ToSqlWhere(c)
	if err != nil {
		return sql4prepare, values, err
	}
	values = append(values, binds...)

	sql4prepare = NamedSprintf("UPDATE :tables SET :updates :wheres", tables, strings.Join(updates, ", "), wheres)

	return
}

func (b Builder) toSqlSqlDelete(c *dbgo2.Context) (sql4prepare string, values []any, err error) {
	var tables string
	tables, _, err = b.ToSqlTable(c)
	if err != nil {
		return
	}
	wheres, binds, err := b.ToSqlWhere(c)
	if err != nil {
		return sql4prepare, values, err
	}
	values = append(values, binds...)
	sql4prepare = NamedSprintf("DELETE FROM :tables :wheres", tables, wheres)
	return
}
