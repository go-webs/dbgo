package dbgo

import (
	"cmp"
	"fmt"
	"gitub.com/go-webs/dbgo/builder"
	"gitub.com/go-webs/dbgo/iface"
	"reflect"
	"regexp"
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
	argsStr := Map[T, []T, string](args, func(s T) string {
		return reflect.ValueOf(s).String()
	})
	return iface.TypeRaw(strings.Join(argsStr, ","))
}

func Map[Data any, Datas ~[]Data, Result any](datas Datas, mapper func(Data) Result) []Result {
	results := make([]Result, 0, len(datas))
	for _, data := range datas {
		results = append(results, mapper(data))
	}
	return results
}

//	func TableAs(table any, as string) []any {
//		return []any{table, as}
//	}
func TableAs(table any, as string) *builder.TableBuilder {
	return builder.NewTableBuilder("").Table(table, as)
}

func NamedSprintf(format string, a ...any) string {
	str := regexp.MustCompile(`:\w+`).ReplaceAllString(format, "%s")
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(fmt.Sprintf(str, a...), " "))
}
