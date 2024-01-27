package builder

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"slices"
	"strings"
)

const TAGIGNORE = "-"

type Bindery struct {
	PrimaryKey      string
	PrimaryKeyValue any

	FieldsStruct []string
	FieldsTag    []string

	Datas []map[string]any
}
type BindBuilder struct {
	*Bindery
}

func NewBindBuilder() *BindBuilder {
	return &BindBuilder{}
}

func (b *BindBuilder) BuildFieldsQuery(rft reflect.Type) (err error) {
	b.Bindery = &Bindery{}
	//rfv := reflect.Indirect(reflect.ValueOf(b.bind))
	switch rft.Kind() {
	case reflect.Struct:
		//rft := rfv.Type()
		for i := 0; i < rft.NumField(); i++ {
			typeField := rft.Field(i)
			if typeField.Anonymous {
				continue
			}
			tags := typeField.Tag.Get("db")
			var tag string
			if tags == "" {
				tag = strings.ToLower(typeField.Name)
			} else {
				tagSplit := strings.Split(tags, ",")
				if len(tagSplit) > 1 {
					b.PrimaryKey = tagSplit[1]
				}
				tag = tagSplit[0]
				if tag == TAGIGNORE || typeField.Name == "TableName" {
					continue
				}
			}
			// query
			b.FieldsStruct = append(b.FieldsStruct, typeField.Name)
			b.FieldsTag = append(b.FieldsTag, tag)
		}
	case reflect.Slice:
		if rft.Elem().Kind() == reflect.Struct {
			return b.BuildFieldsQuery(rft.Elem())
		}
	default:
		err = errors.New("no support bind object")
	}
	return
}
func (b *BindBuilder) BuildFieldsExecute(data any, mustFields ...string) (err error) {
	b.Bindery = &Bindery{}
	rfv := reflect.Indirect(reflect.ValueOf(data))
	switch rfv.Kind() {
	case reflect.Struct:
		b.buildFieldsExecuteStruct(rfv, mustFields...)
	case reflect.Slice:
		if rfv.Len() == 0 {
			err = errors.New("data length 0")
		} else {
			if rfv.Type().Elem().Kind() == reflect.Struct {
				for i := 0; i < rfv.Len(); i++ {
					b.buildFieldsExecuteStruct(rfv.Index(i), mustFields...)
				}
			}
		}
	default:
		err = errors.New("no support bind object")
	}
	return
}
func (b *BindBuilder) buildFieldsExecuteStruct(rfv reflect.Value, mustFields ...string) error {
	entry := map[string]any{}
	rft := rfv.Type()
	for i := 0; i < rft.NumField(); i++ {
		typeField := rft.Field(i)
		if typeField.Anonymous {
			continue
		}
		tags := typeField.Tag.Get("db")
		var tag string
		var isPk bool
		if tags == "" {
			tag = strings.ToLower(typeField.Name)
		} else {
			// primary key
			tagSplit := strings.Split(tags, ",")
			if len(tagSplit) > 1 && tagSplit[1] == "pk" {
				b.PrimaryKey = tagSplit[0]
				b.PrimaryKeyValue = rfv.Field(i).Interface()
				isPk = true
			}
			tag = tagSplit[0]
			if tag == TAGIGNORE {
				continue
			}
		}
		// insert/update
		if (rfv.Field(i).Kind() == reflect.Ptr && rfv.Field(i).IsNil()) || (rfv.Field(i).IsZero() && !slices.Contains(mustFields, tag)) {
			continue
		}
		b.FieldsStruct = append(b.FieldsStruct, typeField.Name)
		b.FieldsTag = append(b.FieldsTag, tag)
		if isPk {
			continue
		}
		var rfvVal = rfv.Field(i).Interface()
		if v, ok := rfvVal.(driver.Valuer); ok {
			value, err := v.Value()
			if err != nil {
				return err
			}
			entry[tag] = value
		} else {
			entry[tag] = rfvVal
		}
	}
	b.Datas = append(b.Datas, entry)
	return nil
}
