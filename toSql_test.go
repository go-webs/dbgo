package dbgo2

import "testing"

func db() Database {
	return Open(nil).NewDatabase()
}
func TestDatabase_ToSql(t *testing.T) {
	prepare, values, err := db().Table("a").Select("b").Where("c", 1).OrderBy("id").Limit(10).Page(2).ToSql()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(prepare)
	t.Log(values)
}
