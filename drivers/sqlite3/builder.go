package sqlite3

import (
	"go-webs/dbgo2/drivers"
	"go-webs/dbgo2/drivers/mysql"
)

const DriverName = "sqlite3"

type Driver struct {
	mysql.Driver
}

func init() {
	drivers.Register(DriverName, &Driver{})
}
