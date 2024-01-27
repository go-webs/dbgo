package dbgo

import (
	"github.com/go-webs/dbgo/iface"
	"testing"
)

func TestDatabase_InnerJoin(t *testing.T) {
	var build = db().Table("users", "a").Join(TableAs("card", "b"), "a.id", "=", "b.user_id")
	fields, _, _ := build.BuildJoin()
	var expect = "INNER JOIN `test_card` b ON a.id = b.user_id"
	assertsEqual(t, expect, fields)
	expect = "SELECT * FROM `test_users` a INNER JOIN `test_card` b ON a.id = b.user_id"
	assertsEqual(t, expect, build.ToSqlOnly())
}

func TestDatabase_LeftJoin(t *testing.T) {
	var build = db().Table("users", "a").
		LeftJoin(TableAs("card", "b"), "a.id", "=", "b.user_id")
	fields, _, _ := build.BuildJoin()
	var expect = "LEFT JOIN `test_card` b ON a.id = b.user_id"
	assertsEqual(t, expect, fields)
	expect = "SELECT * FROM `test_users` a LEFT JOIN `test_card` b ON a.id = b.user_id"
	assertsEqual(t, expect, build.ToSqlOnly())
}

func TestDatabase_RightJoin(t *testing.T) {
	var build = db().Table("users", "a").
		RightJoin(TableAs("card", "b"), "a.id", "=", "b.user_id")
	fields, _, _ := build.BuildJoin()
	var expect = "RIGHT JOIN `test_card` b ON a.id = b.user_id"
	assertsEqual(t, expect, fields)
	expect = "SELECT * FROM `test_users` a RIGHT JOIN `test_card` b ON a.id = b.user_id"
	assertsEqual(t, expect, build.ToSqlOnly())
}

func TestDatabase_CrossJoin(t *testing.T) {
	var build = db().Table("users", "a").
		CrossJoin("card")
	fields, _, _ := build.BuildJoin()
	var expect = "CROSS JOIN `test_card`"
	assertsEqual(t, expect, fields)
	expect = "SELECT * FROM `test_users` a CROSS JOIN `test_card`"
	assertsEqual(t, expect, build.ToSqlOnly())
}

func TestDatabase_Join(t *testing.T) {
	var build = db().Table("users", "a").
		Join(TableAs("card", "b"), func(joins iface.JoinClause) {
			joins.On("a.id", "=", "b.user_id").OrOn("a.age", "=", "b.age")
		}).
		Join(TableAs("address", "c"), func(joins iface.JoinClause) {
			joins.On("a.id", "=", "c.user_id").OrOn("b.age", "=", "c.age")
		})
	fields, _, _ := build.BuildJoin()
	var expect = "JOIN `test_card` b ON (a.id = b.user_id OR a.age = b.age) JOIN `test_address` c ON (a.id = c.user_id OR b.age = c.age)"
	assertsEqual(t, expect, fields)
	expect = "SELECT * FROM `test_users` a JOIN `test_card` b ON (a.id = b.user_id OR a.age = b.age) JOIN `test_address` c ON (a.id = c.user_id OR b.age = c.age)"
	assertsEqual(t, expect, build.ToSqlOnly())
}
func TestDatabase_Union(t *testing.T) {
	var build = db().Table("users").
		Union(db().Table("card", "b"))
	fields, _, _ := build.BuildJoin()
	var expect = "UNION ALL (SELECT * FROM `test_card` b)"
	assertsEqual(t, expect, fields)
	expect = "SELECT * FROM `test_users` UNION ALL (SELECT * FROM `test_card` b)"
	assertsEqual(t, expect, build.ToSqlOnly())
}
