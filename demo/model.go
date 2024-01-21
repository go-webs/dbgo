package demo

import "time"

type Users struct {
	Id        int       `db:"id,primaryKey"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Title     string    `db:"title"`
	Active    bool      `db:"active"`
	Votes     int       `db:"votes"`
	Balance   float64   `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
}

func (Users) TableName() string {
	return "users"
}

type Cards struct {
	Users
	Id   int    `db:"id,primaryKey"`
	Name string `db:"name"`
}

func (Cards) TableName() string {
	return "cards"
}
