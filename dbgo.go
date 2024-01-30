package dbgo2

import (
	"database/sql"
	"go-webs/dbgo2/util"
)

type DbGo struct {
	Cluster *ConfigCluster
	master  []*sql.DB
	slave   []*sql.DB
	driver  string
	prefix  string

	//SqlLogs        []string
	//enableQueryLog bool
	//Error          error
	handlers []HandlerFunc
}

type HandlerFunc func(*Context)

func (dg *DbGo) Use(h ...HandlerFunc) *DbGo {
	dg.handlers = append(dg.handlers, h...)
	return dg
}

// Open db
// examples
//
//	Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true")
//	Open(&ConfigCluster{...})
func Open(conf ...any) *DbGo {
	var dg = DbGo{}
	switch len(conf) {
	case 1:
		if single, ok := conf[0].(*Config); ok {
			dg.driver = single.Driver
			dg.prefix = single.Prefix
			dg.Cluster = &ConfigCluster{WriteConf: []Config{*single}}
			if single == nil { // build sql only
				//dg.ConfigCluster = &ConfigCluster{Prefix: "test_"}
				return &dg
			}
			dg.master, dg.slave = dg.Cluster.init()
		} else if cluster, ok := conf[0].(*ConfigCluster); ok {
			dg.driver = cluster.WriteConf[0].Driver
			dg.prefix = cluster.WriteConf[0].Prefix
			dg.Cluster = cluster
			if cluster == nil { // build sql only
				//dg.ConfigCluster = &ConfigCluster{Prefix: "test_"}
				return &dg
			}
			dg.master, dg.slave = dg.Cluster.init()
		}
	case 2:
		dg.driver = conf[0].(string)
		db, err := sql.Open(dg.driver, conf[1].(string))
		if err != nil {
			panic(err.Error())
		}
		dg.master = append(dg.master, db)
	default:
		panic("config must be *dbgo.ConfigCluster or sql.Open() origin params")
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
func (dg *DbGo) Driver() IDriver {
	return GetDriver(dg.driver)
}

//func (dg *DbGo) NewDB() *Database {
//	return newDatabase(dg)
//}

//func (dg *DbGo) queryRow(db *sql.DB, query string, args ...any) *sql.Row {
//	prepare, err := db.Prepare(query)
//	return dg.SlaveDB().QueryRow(query, args...)
//}
//func (dg *DbGo) QueryRow(query string, args ...any) *sql.Row {
//	return dg.SlaveDB().QueryRow(query, args...)
//}

//func (dg *DbGo) QueryRow(query string, args ...any) *sql.Row {
//	return dg.SlaveDB().QueryRow(query, args...)
//}
//func (dg *DbGo) Query(query string, args ...any) (*sql.Rows, error) {
//	return dg.SlaveDB().Query(query, args...)
//}
//func (dg *DbGo) Exec(query string, args ...any) (sql.Result, error) {
//	return dg.MasterDB().Exec(query, args...)
//}
//func (dg *DbGo) Begin() (*sql.Tx, error) {
//	return dg.MasterDB().Begin()
//}
//func (dg *DbGo) Trans(closer ...func(*sql.Tx) error) error {
//	var tx, err = dg.MasterDB().Begin()
//	if err != nil {
//		return err
//	}
//	for _, v := range closer {
//		err = v(tx)
//		if err != nil {
//			return tx.Rollback()
//		}
//	}
//	return tx.Commit()
//}

func (dg *DbGo) NewDatabase() Database {
	return NewDatabase(dg)
}

func (dg *DbGo) NewSession() *Session {
	master := dg.MasterDB()
	slave := dg.SlaveDB()
	return NewSession(master, slave)
}

//func (dg *DbGo) Transaction() error {
//	tx, err := dg.Begin()
//	if err != nil {
//		return err
//	}
//}

//func (dg *DbGo) EnableQueryLog(b bool) {
//	dg.enableQueryLog = b
//}
//func (dg *DbGo) recordSqlLog(queryStr string, values ...interface{}) {
//	logrus.Debug("record sql log: "+queryStr, values)
//	if dg.enableQueryLog {
//		dg.SqlLogs = append(dg.SqlLogs, fmt.Sprintf("%s, %v", queryStr, values))
//	}
//}
//func (dg *DbGo) LastSql() (last string) {
//	if len(dg.SqlLogs) > 0 {
//		last = dg.SqlLogs[len(dg.SqlLogs)-1]
//	}
//	return
//}
