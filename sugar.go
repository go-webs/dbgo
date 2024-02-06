package dbgo2

func (db Database) MustGet(columns ...string) (res []map[string]any) {
	res, db.Err = db.Get(columns...)
	return
}
