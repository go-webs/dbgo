package dbgo

import (
	"fmt"
)

func (db Database) HasOne(tab2 any, tab1col, tab2col string) (res map[string]any, subKey string, err error) {
	res, err = db.First()
	if err != nil {
		return
	}
	var result map[string]any
	result, err = newDatabase(db.DbGo).Table(tab2).Where(tab2col, res[tab1col]).First()
	if err != nil {
		return
	}
	subKey = fmt.Sprintf("sub_%s", db.TableBuilder.BuildTableNameOnly(tab2))
	if len(result) > 0 {
		res[subKey] = result
	} else {
		res[subKey] = nil
	}
	return
}
func (db Database) HasMany(tab2 any, tab1col, tab2col string) (res map[string]any, subKey string, err error) {
	res, err = db.First()
	if err != nil {
		return
	}
	var result []map[string]any
	result, err = newDatabase(db.DbGo).Table(tab2).Where(tab2col, res[tab1col]).Get()
	if err != nil {
		return
	}
	subKey = fmt.Sprintf("sub_%s", db.TableBuilder.BuildTableNameOnly(tab2))
	if len(result) > 0 {
		res[subKey] = result
	} else {
		res[subKey] = []any{}
	}
	return
}

type HasManyThroughProps struct {
	Local               any
	Through             any
	Children            any
	Local2ThroughKey    []string
	Through2ChildrenKey []string
}

func (db Database) HasManyThrough(Local, Through, Children any, Local2ThroughKey, Through2ChildrenKey []string) (res map[string]any, subKey string, err error) {
	res, err = newDatabase(db.DbGo).Table(Local).First()
	if err != nil {
		return
	}
	var result []map[string]any
	result, err = newDatabase(db.DbGo).
		Select("c.*").
		Table(Local, "a").
		LeftJoin(Through, "b", "a."+Local2ThroughKey[0], "b."+Local2ThroughKey[1]).
		LeftJoin(Children, "b."+Through2ChildrenKey[0], "=", "c."+Through2ChildrenKey[1]).
		Get()
	if err != nil {
		return
	}
	subKey = fmt.Sprintf("sub_%s", db.TableBuilder.BuildTableNameOnly(Children))
	if len(result) > 0 {
		res[subKey] = result
	} else {
		res[subKey] = []any{}
	}
	//return res, subKey
	return
}

func (db Database) HasOneThrough(Local, Through, Children any, Local2ThroughKey, Through2ChildrenKey []string) (res map[string]any, subKey string, err error) {
	res, err = newDatabase(db.DbGo).Table(Local).First()
	if err != nil {
		return
	}
	var result map[string]any
	result, err = newDatabase(db.DbGo).
		Select("c.*").
		Table(Local, "a").
		LeftJoin(Through, "b", "a."+Local2ThroughKey[0], "b."+Local2ThroughKey[1]).
		LeftJoin(Children, "b."+Through2ChildrenKey[0], "=", "c."+Through2ChildrenKey[1]).
		First()
	if err != nil {
		return
	}
	subKey = fmt.Sprintf("sub_%s", db.TableBuilder.BuildTableNameOnly(Children))
	if len(result) > 0 {
		res[subKey] = result
	} else {
		res[subKey] = nil
	}
	return
}
