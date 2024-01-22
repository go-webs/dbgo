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
	fields, _, _ := build.BuildJoin()
	var expect = "LIMIT ? OFFSET ?"
	assertsEqual(t, expect, fields)
	expect = "SELECT * FROM `test_users` a GROUP BY `a`,`b` HAVING `id` IS NOT NULL ORDER BY id ASC LIMIT ? OFFSET ?"
	assertsEqual(t, expect, build.ToSqlOnly())
}

func TestDatabase_BuildSqlInsert(t *testing.T) {
	segment, values, err := db().Table("users").BuildSqlInsert(map[string]any{
		"id":    1,
		"name":  "john",
		"email": "a@a.com",
	})
	if err != nil {
		t.Errorf("TestDatabase_BuildSqlInsert error:%s", err)
	}
	t.Log(segment, values)

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
	t.Log(segment, values)
}
