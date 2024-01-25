package dbgo

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"strings"
	"testing"
)

type Users struct {
	Id        int64         `db:"id,primaryKey" json:"id,omitempty"`
	Name      string        `db:"name" json:"name,omitempty"`
	Email     string        `db:"email" json:"email,omitempty"`
	Votes     sql.NullInt64 `db:"votes" json:"votes,omitempty"`
	CreatedAt sql.NullTime  `db:"created_at" json:"created_at,omitempty"`

	TableName string `db:"users" json:"-"`
}

func db() *Database {
	var dbg = Open(&Cluster{Prefix: "test_"})
	dbg.EnableQueryLog(true)
	return dbg.NewDB()
}

var dbg2 = Open("mysql", "root:Qx233233!@tcp(rm-bp1149oa09n39n236jo.mysql.rds.aliyuncs.com:3306)/game?charset=utf8mb4&parseTime=true")

func db2() *Database {
	dbg2.EnableQueryLog(true)
	return dbg2.NewDB()
}

func assertsEqual(t *testing.T, expect, real any) {
	marshal, err := json.Marshal(expect)
	assertsError(t, err)
	bytes, err := json.Marshal(real)
	assertsError(t, err)
	if string(marshal) != string(bytes) {
		methodName, file, line := getCallerInfo()
		t.Errorf("[%s] Error\n\t Trace - %s:%v\n\tExpect - %T %s\n\t   Got - %T %s\n------------------------------------------------------", methodName, file, line, expect, marshal, real, bytes)
	}
}

func assertsError(t *testing.T, err error) {
	if err != nil {
		methodName, file, line := getCallerInfo()
		t.Errorf("[%s] Error\n\t Trace - %s:%v\n\t%s\n------------------------------------------------------", methodName, file, line, err.Error())
	}
}

func getCallerInfo() (string, string, int) {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var i int
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		i++
		if i == 1 {
			continue
		}
		//fmt.Printf("Method: %s\nFile: %s\nLine: %d\n\n", frame.Function, frame.File, frame.Line)
		lastDotIndex := strings.LastIndex(frame.Function, ".")
		methodName := frame.Function[lastDotIndex+1:]
		//t.Logf("[%s] errors on file:line: \n\t\t -> %s:%v\n", methodName, frame.File, frame.Line)
		if i == 2 {
			return methodName, frame.File, frame.Line
		}
		//break
	}
	return "", "", 0
}

func JsonLog(t *testing.T, data any) {
	marshal, err := json.Marshal(data)
	assertsError(t, err)
	t.Logf("json data: %s", marshal)
}
