package builder

const TAGIGNORE = "-"

type fieldStruct struct {
	FieldsStruct []string
	FieldsTag    []string
}
type BindBuilder struct {
	bind         any
	PrimaryKey   string
	FieldsQuery  fieldStruct
	FieldsInsert fieldStruct
	FieldsUpdate fieldStruct
	FieldsDelete fieldStruct
	Type         int8 // 0:struct, 1:struct slice, 2:map, 3:map slice
}

func NewBindBuilder(bind any) *BindBuilder {
	return &BindBuilder{bind: bind}
}

func (b *BindBuilder) BuildFields() {
	//rfv := reflect.Indirect(reflect.ValueOf(b.bind))
	//switch rfv.Kind() {
	//case reflect.Struct:
	//	rft := rfv.Type()
	//	for i := 0; i < rft.NumField(); i++ {
	//		typeField := rft.Field(i)
	//		if typeField.Anonymous {
	//			continue
	//		}
	//		tags := typeField.Tag.Get("db")
	//		var tag string
	//		if tags == "" {
	//			tag = strings.ToLower(typeField.Name)
	//		} else {
	//			tagSplit := strings.Split(tags, ",")
	//			if len(tagSplit) > 1 {
	//				b.PrimaryKey = tagSplit[1]
	//			}
	//			tag = tagSplit[0]
	//			if tag == TAGIGNORE {
	//				continue
	//			}
	//		}
	//		// query
	//		b.FieldsQuery.FieldsStruct = append(b.FieldsQuery.FieldsStruct, typeField.Name)
	//		b.FieldsQuery.FieldsTag = append(b.FieldsQuery.FieldsTag, tag)
	//		// insert
	//		if (rfv.Field(i).Kind() == reflect.Ptr && rfv.Field(i).IsNil()) || (rfv.Field(i).IsZero() && !slices.Contains(mustFields, column)) {
	//			continue
	//		}
	//		b.FieldsInsert.FieldsStruct = append(b.FieldsInsert.FieldsStruct, typeField.Name)
	//		b.FieldsInsert.FieldsTag = append(b.FieldsInsert.FieldsTag, tag)
	//	}
	//case reflect.Slice:
	//	if rfv.Type().Elem().Kind() == reflect.Struct {
	//
	//	}
	//default:
	//
	//}
}
