package dbgo

import (
	"testing"
)

func TestDatabase_Insert(t *testing.T) {
	var u = []Users{
		//{Email: "a112121@a.com", Votes: 11},
		//{Email: "a112132@a.com", Votes: 11},
	}
	var d = db()
	err2 := d.BuildFieldsExecute(&u)
	assertsError(t, err2)
	t.Logf("%+v", *d.BindBuilder.Bindery)

	var d2 = db2()
	rows, err := d2.Insert(&u)

	t.Log(d2.SqlLogs)

	assertsError(t, err)
	assertsEqual(t, int64(1), rows)
}

func TestDatabase_To(t *testing.T) {
	var u Users
	var d = db2()
	err := d.To(&u)
	assertsError(t, err)
	assertsEqual(t, int64(1), u.Id)
	t.Logf("%+v", u)

	var u2 = Users{Id: 2}
	err = db2().To(&u2)
	assertsError(t, err)
	assertsEqual(t, 2, u2.Id)
}
