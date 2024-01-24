package dbgo

import (
	"cmp"
	"database/sql"
	"fmt"
	"gitub.com/go-webs/dbgo/builder"
	"gitub.com/go-webs/dbgo/iface"
	"gitub.com/go-webs/dbgo/util"
	"reflect"
	"strings"
)

type DbGo struct {
	*Cluster
	master  []*sql.DB
	slave   []*sql.DB
	SqlLogs []string

	enableQueryLog bool
	Error          error
}

// Open db
// examples
//
//	Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true")
//	Open(&Cluster{...})
func Open(conf ...any) *DbGo {
	var dg = DbGo{Cluster: &Cluster{}}
	switch len(conf) {
	case 1:
		if cluster, ok := conf[0].(*Cluster); ok {
			dg.Cluster = cluster
			if cluster == nil { // build sql only
				//dg.Cluster = &Cluster{Prefix: "test_"}
				return &dg
			}
			dg.master, dg.slave = cluster.init()
		}
	case 2:
		db, err := sql.Open(conf[0].(string), conf[1].(string))
		if err != nil {
			panic(err.Error())
		}
		dg.master = append(dg.master, db)
		//return Open(&Cluster{Master: []Config{{Host: conf[1].(string)}}, Driver: conf[0].(string)})
	default:
		panic("config must be *dbgo.Cluster or sql.Open() origin params")
	}
	return &dg
}

func (dg *DbGo) MasterDB() *sql.DB {
	return dg.master[util.GetRandomInt(len(dg.master))]
}
func (dg *DbGo) SlaveDB() *sql.DB {
	return dg.slave[util.GetRandomInt(len(dg.slave))]
}
func (dg *DbGo) NewDB() *Database {
	return newDatabase(dg)
}
func (dg *DbGo) EnableQueryLog(b bool) {
	dg.enableQueryLog = b
}
func (dg *DbGo) recordSqlLog(queryStr string, values ...interface{}) {
	if dg.enableQueryLog {
		dg.SqlLogs = append(dg.SqlLogs, fmt.Sprintf("%s, %v", queryStr, values))
	}
}
func (dg *DbGo) LastSql() (last string) {
	if len(dg.SqlLogs) > 0 {
		last = dg.SqlLogs[len(dg.SqlLogs)-1]
	}
	return
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
