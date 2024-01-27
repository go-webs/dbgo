package dbgo

import (
	"cmp"
	"github.com/go-webs/dbgo/builder"
	"github.com/go-webs/dbgo/iface"
	"github.com/go-webs/dbgo/util"
	"reflect"
	"strings"
)

func Raw[T cmp.Ordered](args ...T) iface.TypeRaw {
	argsStr := util.Map[T, []T, string](args, func(s T) string {
		return reflect.ValueOf(s).String()
	})
	return iface.TypeRaw(strings.Join(argsStr, ","))
}

//	func TableAs(table any, as string) []any {
//		return []any{table, as}
//	}
func TableAs(table any, as string) *builder.TableBuilder {
	return builder.NewTableBuilder("").Table(table, as)
}
