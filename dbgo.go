package dbgo

import (
	"cmp"
	"gitub.com/go-webs/dbgo/builder"
	"gitub.com/go-webs/dbgo/iface"
	"gitub.com/go-webs/dbgo/util"
	"reflect"
	"strings"
)

type DbGo struct {
	*Cluster
	Err error
}

func Open(conf *Cluster) *DbGo {
	if conf == nil { // build sql only
		return &DbGo{Cluster: &Cluster{Prefix: "test_"}}
	}
	//todo
	return &DbGo{Cluster: conf}
}

func (dg *DbGo) NewDB() *Database {
	return NewDB(dg)
}
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
