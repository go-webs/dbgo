package doc

import "database/sql"

// Column 表示SELECT语句中的列信息。
type Column struct {
	Name  string
	Alias string // 可选别名
	IsRaw bool   // 是否是原生SQL片段
}

// SelectClause 存储SELECT子句相关信息。
type SelectClause struct {
	Columns  []Column
	Distinct bool
}

// Condition 用于表示WHERE或HAVING子句中的单个条件。
type Condition struct {
	Column    string
	Value     interface{}
	Operator  string
	LogicalOp string        // "AND" 或 "OR"
	Not       bool          // 是否取反操作
	SubQuery  *QueryBuilder // 若条件是一个子查询，则存储该子查询
}

// WhereClause 存储所有WHERE条件。
type WhereClause struct {
	Conditions []Condition
}

// HavingClause 类似于WhereClause，但应用于HAVING子句。
type HavingClause struct {
	Conditions []Condition
}

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

// QueryBuilder 是核心接口，包含所有方法用于构建查询。
type QueryBuilder struct {
	DB           DBInterface // 数据库连接接口
	Selects      SelectClause
	From         string
	Joins        []JoinClause
	Wheres       WhereClause
	GroupBys     []string
	Havings      HavingClause
	OrderBys     OrderByClause
	LimitOffset  LimitOffsetClause
	UnionClauses []*QueryBuilder // 原始SQL UNION/UNION ALL子句
}

// DBInterface 定义了与数据库交互的基本接口。
type DBInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

// Paginator 是用于分页查询结果的结构体，包含当前页数据及分页信息。
type Paginator struct {
	Items       []any
	Total       int64
	CurrentPage int
	PerPage     int
	LastPage    int
}

// NamedQuery 结构体，用于存储命名查询及其参数
type NamedQuery struct {
	Name     string
	SQL      string
	Bindings []interface{}
}

type SubQueryBuilder struct {
	Parent  *QueryBuilder
	Builder *QueryBuilder
	//...
}

type WhenClause struct {
	Condition interface{}
	Result    interface{}
}
type CaseStatement struct {
	Builder *QueryBuilder
	Whens   []*WhenClause
	Else    interface{}
	Column  string // 可选的列名，用于定义CASE语句作用于哪一列
}

type WindowFunctionBuilder struct {
	Builder      *QueryBuilder
	FunctionName string
	PartitionBy  []string
	OrderBy      []OrderByItem
	Frame        *WindowFrame
}

type OrderByItem struct {
	Column    string
	Direction string // "asc" 或 "desc"
}

type Expression struct {
	SQL        string
	BindValues []interface{}
}

type PagingResult struct {
	Data        []any // 当前页的数据
	TotalCount  int64 // 数据总数
	PageCount   int   // 总页数
	CurrentPage int   // 当前页码
	PerPage     int   // 每页数量
}
