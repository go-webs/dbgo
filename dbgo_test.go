package dbgo

import (
	"cmp"
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

func assertsEqual[T cmp.Ordered](t *testing.T, expect, real T) {
	if expect != real {
		t.Errorf("not equal, expect: %v\n but got: %v", expect, real)
	}
}
