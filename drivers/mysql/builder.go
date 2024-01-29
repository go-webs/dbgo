package mysql

import (
	"go-webs/dbgo2/drivers"
)

const DriverName = "mysql"

type Driver struct {
	BB string
}

func init() {
	drivers.Register(DriverName, &Driver{})
}


func (Driver) ToSql(ctx *drivers.Context) (sql4prepare string, binds []any, err error) {
	return
}
func (Driver) AA() {

}
