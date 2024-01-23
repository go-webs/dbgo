package dbgo

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// 'mysql' => [
//    'read' => [
//        'host' => [
//            '192.168.1.1',
//            '196.168.1.2',
//        ],
//    ],
//    'write' => [
//        'host' => [
//            '196.168.1.3',
//        ],
//    ],
//    'sticky' => true,
//    'driver' => 'mysql',
//    'database' => 'database',
//    'username' => 'root',
//    'password' => '',
//    'charset' => 'utf8mb4',
//    'collation' => 'utf8mb4_unicode_ci',
//    'prefix' => '',
//],

//type Config struct {
//	Host      []string
//	Sticky    bool   `toml:"sticky"`
//	Driver    string `toml:"driver"`
//	Database  string `toml:"database"`
//	Username  string `toml:"username"`
//	Password  string `toml:"password"`
//	Charset   string `toml:"charset"`
//	Collation string `toml:"collation"`
//	Prefix    string `toml:"prefix"`
//}

type Config struct {
	Host            string
	Port            int
	MaxIdleConns    *int
	MaxOpenConns    *int
	ConnMaxLifetime *time.Duration
	ConnMaxIdleTime *time.Duration
}

type Cluster struct {
	Master    []Config
	Slave     []Config
	Sticky    bool   `toml:"sticky"`
	Driver    string `toml:"driver"`
	Database  string `toml:"database"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	Charset   string `toml:"charset"`
	Collation string `toml:"collation"`
	Prefix    string `toml:"prefix"`
	ParseTime bool   `toml:"parseTime"`
}

func (c Cluster) init() (master []*sql.DB, slave []*sql.DB) {
	var optionsArr []string
	if c.Charset != "" {
		optionsArr = append(optionsArr, fmt.Sprintf("charset=%s", c.Charset))
	}
	if c.Collation != "" {
		optionsArr = append(optionsArr, fmt.Sprintf("collation=%s", c.Collation))
	}
	if c.ParseTime {
		optionsArr = append(optionsArr, fmt.Sprintf("parseTime=%v", c.Charset))
	}

	var options string
	if len(optionsArr) > 0 {
		options = fmt.Sprintf("?%s", strings.Join(optionsArr, "&"))
	}
	if len(c.Master) > 0 {
		for _, v := range c.Master {
			master = append(master, c.initDB(v, options))
		}
	}
	if len(c.Slave) > 0 {
		for _, v := range c.Master {
			master = append(master, c.initDB(v, options))
		}
	}
	return
}
func (c Cluster) initDB(v Config, options string) *sql.DB {
	if v.Port == 0 {
		v.Port = 3306
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s%s", c.Username, c.Password, v.Host, v.Port, c.Database, options)
	db, err := sql.Open(c.Driver, dsn)
	if err != nil {
		panic(err.Error())
	}

	if v.MaxIdleConns != nil {
		db.SetMaxIdleConns(*v.MaxIdleConns)
	}
	if v.MaxOpenConns != nil {
		db.SetMaxOpenConns(*v.MaxOpenConns)
	}
	if v.ConnMaxLifetime != nil {
		db.SetConnMaxLifetime(*v.ConnMaxLifetime)
	}
	if v.ConnMaxIdleTime != nil {
		db.SetConnMaxIdleTime(*v.ConnMaxIdleTime)
	}

	return db
}
