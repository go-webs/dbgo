package dbgo

import (
	"database/sql"
	"math"
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
	handlers HandlersChain
}

type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc

const abortIndex int8 = math.MaxInt8 >> 1

//type handlers struct {
//	handlers HandlersChain
//	index    int8
//}
//
//func (c *handlers) Next() {
//	c.index++
//	for c.index < int8(len(c.handlers)) {
//		c.handlers[c.index](c)
//		c.index++
//	}
//}

func (dg *DbGo) Use(h ...HandlerFunc) *DbGo {
	dg.handlers = append(dg.handlers, h...)
	return dg
}

//func (dg *DbGo) Use(h ...HandlerFunc) *DbGo {
//}

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
		} else {
			dg.driver = "mysql" // for toSql test
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
	if len(dg.master) == 0 {
		return nil
	}
	return dg.master[GetRandomInt(len(dg.master))]
}
func (dg *DbGo) SlaveDB() *sql.DB {
	if len(dg.slave) == 0 {
		return dg.MasterDB()
	}
	return dg.slave[GetRandomInt(len(dg.slave))]
}
func (dg *DbGo) Driver() IDriver {
	return GetDriver(dg.driver)
}

func (dg *DbGo) NewDatabase() Database {
	return NewDatabase(dg)
}

func (dg *DbGo) NewSession() *Session {
	return NewSession(dg)
}
