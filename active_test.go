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

	var u2 = Users{Id: 2}
	err = db2().To(&u2)
	assertsError(t, err)
	assertsEqual(t, int64(2), u2.Id)

	var u3 []Users
	err = db2().Limit(2).To(&u3)
	assertsError(t, err)
	assertsEqual(t, int64(1), u3[0].Id)
}

func TestDatabase_First(t *testing.T) {
	res, err := db2().Table(Users{}).First()
	assertsError(t, err)
	assertsEqual(t, int64(1), res["id"])
	t.Log(db2().LastSql())
}
func TestDatabase_Find(t *testing.T) {
	res, err := db2().Table(Users{}).Find(1)
	assertsError(t, err)
	assertsEqual(t, int64(1), res["id"])
	t.Log(db2().LastSql())
}

func TestDatabase_Get(t *testing.T) {
	res, err := db2().Table(Users{}).Limit(2).Get()
	JsonLog(t, res)
	assertsError(t, err)
	assertsEqual(t, int64(1), res[0]["id"])
	t.Log(db2().LastSql())
}

func TestDatabase_Count(t *testing.T) {
	res, err := db2().Table(Users{}).Count()
	JsonLog(t, res)
	assertsError(t, err)
	assertsEqual(t, int64(14), res)
	t.Log(db2().LastSql())
}

func TestDatabase_Value(t *testing.T) {
	res, err := db2().Table(Users{}).Where("id", 1).Value("id")
	assertsError(t, err)
	assertsEqual(t, int64(1), res)
	t.Log(db2().LastSql())
}

func TestDatabase_Pluck(t *testing.T) {
	res, err := db2().Table(Users{}).OrderBy("id").Limit(2).Pluck("id")
	//JsonLog(t, res)
	assertsError(t, err)
	assertsEqual(t, int64(1), res.([]any)[0])
	t.Log(db2().LastSql())
}

func TestDatabase_Max(t *testing.T) {
	res, err := db2().Table(Users{}).Max("id")
	//JsonLog(t, res)
	assertsError(t, err)
	assertsEqual(t, float64(20), res)
	t.Log(db2().LastSql())
}
func TestDatabase_Min(t *testing.T) {
	res, err := db2().Table(Users{}).Min("id")
	//JsonLog(t, res)
	assertsError(t, err)
	assertsEqual(t, float64(1), res)
	t.Log(db2().LastSql())
}
func TestDatabase_Avg(t *testing.T) {
	res, err := db2().Table(Users{}).Avg("id")
	//JsonLog(t, res)
	assertsError(t, err)
	assertsEqual(t, 11.0714, res)
	t.Log(db2().LastSql())
}

func TestDatabase_Exists(t *testing.T) {
	res, err := db2().Table(Users{}).Where("id", 1).Exists()
	//JsonLog(t, res)
	assertsError(t, err)
	assertsEqual(t, true, res)
	t.Log(db2().LastSql())
}
func TestDatabase_DoesntExist(t *testing.T) {
	res, err := db2().Table(Users{}).Where("id", 111).DoesntExist()
	//JsonLog(t, res)
	assertsError(t, err)
	assertsEqual(t, true, res)
	t.Log(db2().LastSql())
}
