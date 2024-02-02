package dbgo2

import (
	"sync"
)

var driverMap = map[string]IDriver{}
var driverLock sync.RWMutex

func Register(driver string, parser IDriver) {
	driverLock.Lock()
	defer driverLock.Unlock()
	driverMap[driver] = parser
}

func GetDriver(driver string) IDriver {
	driverLock.RLock()
	defer driverLock.RUnlock()
	return driverMap[driver]
}
