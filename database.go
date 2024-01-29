package dbgo2

import (
	"go-webs/dbgo2/drivers"
)

type Database struct {
	*DbGo
	*Session

	*drivers.Context
	drivers.IDriver
}

func NewDatabase(dg *DbGo) Database {
	return Database{
		dg,
		dg.NewSession(),
		&drivers.Context{Prefix: dg.prefix},
		drivers.GetDriver(dg.driver),
	}
}

func (db Database) Table() Database {

	return db
}

func (db Database) Where() Database {

	return db
}

func (db Database) Get() Database {
	db.ToSql(db.Context)
	return db
}
