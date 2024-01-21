package dbgo

import "testing"

func TestDatabase_Select(t *testing.T) {
	var build = db().Table("users").Select("id").AddSelect("name as n")
	fields, _ := build.BuildSelect()
	var expect = "`id`,name as n"
	assertsEqual(t, expect, fields)
	expect = "SELECT `id`,name as n FROM `test_users`"
	assertsEqual(t, expect, build.ToSql())
}

func TestDatabase_SelectRaw(t *testing.T) {
	var user Users
	var build = db().Table(&user).Select("id").AddSelect("name as n").SelectRaw("price * ? as price_with_tax", 1.0825)
	fields, binds := build.BuildSelect()
	var expect = "`id`,name as n,price * ? as price_with_tax"
	assertsEqual(t, expect, fields)
	expect = "SELECT `id`,name as n,price * ? as price_with_tax FROM `test_users`"
	assertsEqual(t, expect, build.ToSql())
	assertsEqual(t, 1.0825, binds[0].(float64))
}
