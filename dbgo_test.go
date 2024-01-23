package dbgo

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

type Users struct {
	Id        int       `db:"id,primaryKey"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Title     string    `db:"title"`
	Active    bool      `db:"active"`
	Votes     int       `db:"votes"`
	Balance   float64   `db:"balance"`
	CreatedAt time.Time `db:"created_at"`

	TableName string `db:"test_users"`
}

func db() *Database {
	return Open(&Cluster{Prefix: "test_"}).NewDB()
}

func assertsEqual(t *testing.T, expect, real any) {
	if reflect.ValueOf(expect).String() != reflect.ValueOf(real).String() {
		methodName, file, line := getCallerInfo(t)
		t.Errorf("[%s] Error\n\t Trace - %s:%v\n\tExpect - %+v\n\t   Got - %#v\n------------------------------------------------------", methodName, file, line, expect, real)
	}
}

func getCallerInfo(t *testing.T) (string, string, int) {
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
