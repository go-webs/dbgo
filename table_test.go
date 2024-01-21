package dbgo

import "testing"

func TestDatabase_Table(t *testing.T) {
	var build = db().Table("users")
	var expect = "`test_users`"
	assertsEqual(t, expect, build.BuildTable())
	expect = "SELECT * FROM `test_users`"
	assertsEqual(t, expect, build.ToSql())

	var user Users
	var build2 = db().Table(&user)
	var expect2 = "`test_users`"
	assertsEqual(t, expect2, build2.BuildTable())
	expect2 = "SELECT * FROM `test_users`"
	assertsEqual(t, expect2, build2.ToSql())
}

func TestDatabase_TableAs(t *testing.T) {
	var build = db().Table("users", "a")
	var expect = "`test_users` a"
	assertsEqual(t, expect,
		build.BuildTable())
	expect = "SELECT * FROM `test_users` a"
	assertsEqual(t, expect, build.ToSql())

	var users []Users
	var build2 = db().Table(&users, "a")
	var expect2 = "`test_users` a"
	assertsEqual(t, expect2,
		build2.BuildTable())
	expect2 = "SELECT * FROM `test_users` a"
	assertsEqual(t, expect2, build2.ToSql())
}
