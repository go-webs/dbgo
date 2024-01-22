package dbgo

import "testing"

func TestDatabase_GroupBy(t *testing.T) {
	var build = db().Table("users", "a").GroupBy("a", "b").HavingNotNull("id").HavingRaw("c>1")
	fields, _, _ := build.BuildGroup()
	var expect = "INNER JOIN `test_card` b ON a.id = b.user_id"
	assertsEqual(t, expect, fields)
	expect = "SELECT * FROM `test_users` a INNER JOIN `test_card` b ON a.id = b.user_id"
	assertsEqual(t, expect, build.ToSqlOnly())
}
