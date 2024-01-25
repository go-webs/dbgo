package dbgo

import (
	"gitub.com/go-webs/dbgo/iface"
	"testing"
)

func TestDatabase_BuildSqlQuery(t *testing.T) {
	var build = db().
		Table("users", "a").
		Join(TableAs("card", "b"), "a.id", "=", "b.user_id").
		Distinct().
		Select("a.id", "a.name", "b.no").
		Where("a.votes", ">", db().Table("votes").Select("votes").Limit(1)).
		//Where("b.card_band", "in", func(query *Database) iface.IUnion {
		//	return query.Table("band").Select("name")
		//}).
		Where("b.card_band", "in", db().Table("band").Select("name")).
		Where(func(wh iface.WhereClause) {
			wh.Where("b.status", 1).OrWhere("b.id", ">", 10)
		}).
		GroupBy("a.id").
		Having("a.votes", ">", 10).
		HavingNotNull("b.user_id").OrderByAsc("a.id").Limit(10).Page(2)
	var expect = "SELECT DISTINCT `a`.`id`,`a`.`name`,`b`.`no` FROM `test_users` a INNER JOIN `test_card` b ON `a`.`id` = `b`.`user_id` WHERE a.votes > (SELECT `votes` FROM `test_votes` LIMIT ?) AND b.card_band in (SELECT `name` FROM `test_band`) AND (b.status = ? OR b.id > ?) GROUP BY `a`.`id` HAVING `a`.`votes` > ? AND `b`.`user_id` IS NOT NULL ORDER BY `a`.`id` ASC LIMIT ? OFFSET ?"
	assertsEqual(t, expect, build.ToSqlOnly())
}
func TestDatabase_BuildSqlExists(t *testing.T) {
	var build = db().
		Table("users", "a").
		Join(TableAs("card", "b"), "a.id", "=", "b.user_id").
		Distinct().
		Select("a.id", "a.name", "b.no")
	fields, _, _ := build.BuildSqlExists()
	var expect = "SELECT EXISTS(SELECT DISTINCT `a`.`id`,`a`.`name`,`b`.`no` FROM `test_users` a INNER JOIN `test_card` b ON `a`.`id` = `b`.`user_id`) AS exists"
	assertsEqual(t, expect, fields)
}
func TestDatabase_BuildSqlInsert(t *testing.T) {
	segment, values, err := db().Table("users").BuildSqlInsert(map[string]any{
		"id":    1,
		"name":  "john",
		"email": "a@a.com",
	})
	var expect = "INSERT INTO `test_users` (`email`,`id`,`name`) VALUES (?,?,?)"
	if err != nil {
		t.Errorf("TestDatabase_BuildSqlInsert error:%s", err)
	}
	assertsEqual(t, expect, segment)

	segment, values, err = db().Table("users").BuildSqlInsert([]map[string]any{
		{
			"id":    1,
			"name":  "john",
			"email": "a@a.com",
		},
		{
			"id":    2,
			"name":  "dawn",
			"email": "b@a.com",
		},
	})
	if err != nil {
		t.Errorf("TestDatabase_BuildSqlInsert error:%s", err)
	}
	expect = "INSERT INTO `test_users` (`email`,`id`,`name`) VALUES (?,?,?) (?,?,?)"
	assertsEqual(t, expect, segment)
	for i, v := range []any{"a@a.com", 1, "john", "b@a.com", 2, "dawn"} {
		assertsEqual(t, v, values[i])
	}
}
func TestDatabase_BuildSqlInsertOrIgnore(t *testing.T) {
	var build = db().Table("users")
	fields, _, _ := build.BuildSqlInsertOrIgnore(map[string]any{
		"id":    1,
		"name":  "john",
		"email": "a@a.com",
	})
	var expect = "INSERT IGNORE INTO `test_users` (`email`,`id`,`name`) VALUES (?,?,?)"
	assertsEqual(t, expect, fields)
}
func TestDatabase_BuildSqlUpsert(t *testing.T) {
	var build = db().Table("users")
	fields, _, _ := build.BuildSqlUpsert(map[string]any{
		"id":    1,
		"name":  "john",
		"email": "a@a.com",
	}, []string{}, []string{"name"})
	var expect = "INSERT INTO `test_users` (`email`,`id`,`name`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`)"
	assertsEqual(t, expect, fields)
}
func TestDatabase_BuildSqlInsertUsing(t *testing.T) {
	var build = db().Table("users")
	fields, _, _ := build.BuildSqlInsertUsing([]string{"name", "votes"}, db().Table("users").Select("name", "votes").Limit(3))
	var expect = "INSERT INTO `test_users` (`name`,`votes`) (SELECT `name`,`votes` FROM `test_users` LIMIT ?)"
	assertsEqual(t, expect, fields)
}
func TestDatabase_BuildSqlUpdate(t *testing.T) {
	var build = db().Table("users").Where("id", 1)
	fields, _, _ := build.BuildSqlUpdate(map[string]any{
		"id":    1,
		"name":  "john",
		"email": "a@a.com",
	})
	var expect = "UPDATE `test_users` SET `email` = ?, `id` = ?, `name` = ? WHERE `id` = ?"
	assertsEqual(t, expect, fields)
}
func TestDatabase_BuildSqlDelete(t *testing.T) {
	var build = db().Table("users").Where("name", "John")
	fields, _, _ := build.BuildSqlDelete(2)
	var expect = "DELETE FROM `test_users` WHERE `name` = ? AND `id` = ?"
	assertsEqual(t, expect, fields)
}
func TestDatabase_Increment(t *testing.T) {
	var build = db().Table("users").Where("name", "John")
	fields, _, _ := build.BuildSqlIncrement("age")
	var expect = "UPDATE `test_users` SET + `age` = `age` + 1 WHERE `name` = ?"
	assertsEqual(t, expect, fields)
	fields, _, _ = build.BuildSqlIncrement("age", 2)
	expect = "UPDATE `test_users` SET + `age` = `age` + 2 WHERE `name` = ?"
	assertsEqual(t, expect, fields)
	fields, _, _ = build.BuildSqlIncrement("age", 2, map[string]int{"votes": 10})
	expect = "UPDATE `test_users` SET + `age` = `age` + 2, `votes` = ? WHERE `name` = ?"
	assertsEqual(t, expect, fields)
}
func TestDatabase_IncrementEach(t *testing.T) {
	var build = db().Table("users").Where("name", "John")
	fields, _, _ := build.BuildSqlIncrementEach(map[string]int{"age": 2, "votes": 1})
	var expect = "UPDATE `test_users` SET `age` = `age` + 2, `votes` = `votes` + 1 WHERE `name` = ?"
	assertsEqual(t, expect, fields)
}
