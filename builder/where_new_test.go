package builder

import (
	"github.com/go-webs/dbgo/iface"
	"testing"
)

func TestWhereBuilderNew(t *testing.T) {
	var w = NewWhereBuilderNew()
	w.Where("1=1").Where("a", "in", []int{1, 2}).OrWhere(func(wh iface.WhereClause2) {
		wh.WhereRaw("a>?", 1).OrWhere("b>?", 2)
	}).Where("a", 1).OrWhere("a", ">", 3)
	t.Log(w.BuildWhere())
}
