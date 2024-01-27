package dbgo

import (
	"database/sql"
	"fmt"
	"github.com/go-webs/dbgo/util"
)

type DbGo struct {
	Cluster *Cluster
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
	var dg = DbGo{}
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
		dg.Cluster = &Cluster{Driver: conf[0].(string)}
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
	if len(dg.slave) == 0 {
		return dg.MasterDB()
	}
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
func (dg *DbGo) Ping() error {
	return dg.MasterDB().Ping()
}
