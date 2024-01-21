package dbgo

import "testing"

func TestDatabase_Where(t *testing.T) {
	var build = db().Table("users").Where("a", 1).WhereIn("b", db().Table("card").Select("id").Where("status", 1))
	var expect = "`a` = ? AND `b` IN (SELECT `id` FROM `test_card` WHERE `status` = ?)"
	assertsEqual(t, expect, build.BuildWhereOnly())
	expect = "SELECT * FROM `test_users` WHERE `a` = ? AND `b` IN (SELECT `id` FROM `test_card` WHERE `status` = ?)"
	assertsEqual(t, expect, build.ToSqlOnly())
}
