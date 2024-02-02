package dbgo2

type Context struct {
	SelectClause
	TableClause
	JoinClause
	WhereClause
	HavingClause
	OrderByClause
	LimitOffsetClause
	Groups []string

	Prefix   string
	Queries  string
	Bindings []any
	Err      error
}

type TableClause struct {
	Tables any // table name or struct(slice) or subQuery
	Alias  string
}

// Column 表示SELECT语句中的列信息。
type Column struct {
	Name  string
	Alias string // 可选别名
	IsRaw bool   // 是否是原生SQL片段
	Binds []any  // 绑定数据
}

// SelectClause 存储SELECT子句相关信息。
type SelectClause struct {
	Columns  []Column
	Distinct bool
}

//// Condition 用于表示WHERE或HAVING子句中的单个条件。
//type Condition struct {
//	Column    string
//	Value     interface{}
//	Operator  string   // = > <...
//	LogicalOp string   // "AND" 或 "OR"
//	Not       bool     // 是否取反操作
//	SubQuery  Database // 若条件是一个子查询，则存储该子查询
//}

// JoinClause 描述JOIN操作。
type JoinClause struct {
	Type         string // JOIN类型（INNER, LEFT, RIGHT等）
	Table        string
	FirstColumn  string
	Operator     string
	SecondColumn string
	On           WhereClause
	Alias        string // 别名
}

type OrderByItem struct {
	Column    string
	Direction string // "asc" 或 "desc"
}

// OrderByClause 存储排序信息。
type OrderByClause struct {
	Columns []OrderByItem
}

// LimitOffsetClause 存储LIMIT和OFFSET信息。
type LimitOffsetClause struct {
	Limit  int
	Offset int
	Page   int
}

type Expression struct {
	SQL        string
	BindValues []interface{}
}

// Paginator 是用于分页查询结果的结构体，包含当前页数据及分页信息。
type Paginator struct {
	Items       []any
	Total       int64
	CurrentPage int
	PerPage     int
	LastPage    int
}

// HavingClause 类似于WhereClause，但应用于HAVING子句。
type HavingClause struct {
	*WhereClause
}

// WhereClause 存储所有WHERE条件。
type WhereClause struct {
	Conditions []any
	Err        error
}
type TypeWhereRaw struct {
	LogicalOp string
	Not       bool
	Column    string
	Bindings  []any
}
type TypeWhereNested struct {
	LogicalOp string
	Not       bool
	Column    func(where IWhere)
}
type TypeWhereSubQuery struct {
	LogicalOp string
	Not       bool
	Column    string
	Operator  string
	SubQuery  IBuilder
}
type TypeWhereStandard struct {
	LogicalOp string
	Not       bool
	Column    string
	Operator  string
	Value     any
}
type TypeWhereIn struct {
	LogicalOp string
	Not       bool
	Column    string
	Operator  string
	Value     any
}
type TypeWhereBetween struct {
	LogicalOp string
	Not       bool
	Column    string
	Operator  string
	Value     any
}
