package dbgo2

import (
	"sync"
)

type IDriver interface {
	ToSql(ctx *Context) (sql4prepare string, binds []any, err error)
}

var driverMap map[string]IDriver
var driverLock sync.RWMutex

func Register(driver string, parser IDriver) {
	driverLock.Lock()
	defer driverLock.Unlock()
	driverMap[driver] = parser
}

func GetDriver(driver string) IDriver {
	driverLock.RLock()
	defer driverLock.Unlock()
	return driverMap[driver]
}
