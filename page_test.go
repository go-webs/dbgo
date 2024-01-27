package dbgo

import "testing"

func TestDatabase_Page(t *testing.T) {
	var build = db().Table("users", "a").GroupBy("a", "b").HavingNotNull("id").OrderByAsc("id").Limit(10).Page(2)
	fields, _ := build.BuildPage()
	var expect = "LIMIT ? OFFSET ?"
	assertsEqual(t, expect, fields)
	expect = "SELECT * FROM `test_users` a GROUP BY `a`,`b` HAVING `id` IS NOT NULL ORDER BY `id` ASC LIMIT ? OFFSET ?"
	assertsEqual(t, expect, build.ToSqlOnly())
}
