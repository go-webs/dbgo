package iface

// Record 表示数据库表中的记录，对应于 PHP 中的 Eloquent 模型或查询结果。
type Record interface{}

// QueryBuilder 是用于构建和执行数据库查询的接口。
type QueryBuilder interface {
	// Table 设置要操作的数据库表名。
	//
	// table: 要操作的数据库表名。
	Table(table string) QueryBuilder

	// Select 指定要从数据库中选择的列。
	//
	// columns: 要选择的列名数组或字符串。
	Select(columns ...string) QueryBuilder

	// Where 添加一个"where"条件到查询语句中。
	//
	// column: 列名或者闭包表达式。
	// operator: 操作符，默认为 '='。
	// value: 条件值。
	// boolean: 条件之间的逻辑关系，默认为 'and'。
	Where(column interface{}, operator string, value interface{}, boolean ...string) QueryBuilder

	// WhereBetween 在指定列的值位于给定范围内时添加一个"where"条件。
	//
	// column: 列名。
	// values: 区间范围数组。
	// not: 是否取反，默认为 false。
	WhereBetween(column string, values []interface{}, not ...bool) QueryBuilder

	// OrderBy 对查询结果进行排序。
	//
	// column: 排序的列名。
	// direction: 排序方向（asc 或 desc）。
	OrderBy(column string, direction string) QueryBuilder

	// GroupBy 将查询结果按指定列分组。
	//
	// columns: 分组列名数组。
	GroupBy(columns ...string) QueryBuilder

	// Aggregate 执行聚合函数（如 COUNT, MAX, SUM, AVG, MIN）。
	//
	// function: 聚合函数名称。
	// columns: 聚合的列名数组或字符串。
	Aggregate(function string, columns ...string) (Record, error)

	// Insert 插入新的记录到数据库表中。
	//
	// values: 要插入的数据，通常为键值对组成的 map 或结构体数组。
	Insert(values interface{}) (int64, error)

	// Update 更新符合条件的记录。
	//
	// values: 要更新的数据，通常为键值对组成的 map 或结构体。
	// where: 更新条件（可能是一个或多个Where方法调用链）。
	Update(values interface{}, where QueryBuilder) (int64, error)

	// Delete 删除符合条件的记录。
	//
	// where: 删除条件（可能是一个或多个Where方法调用链）。
	Delete(where QueryBuilder) (int64, error)

	// Get 获取查询结果集。
	//
	// columns: 要获取的列名数组，如果不提供，则获取所有列。
	Get(columns ...string) ([]Record, error)

	// SelectRaw 允许直接在查询中插入原始SQL片段作为选择列。
	//
	// raw: 原始SQL片段。
	SelectRaw(raw string, bindings ...interface{}) QueryBuilder

	// OrWhere 添加一个“或where”条件到查询语句中。
	//
	// column: 列名或者闭包表达式。
	// operator: 操作符，默认为 '='。
	// value: 条件值。
	// boolean: 通常无需指定，因为这个方法默认会使用 'or' 连接条件。
	OrWhere(column interface{}, operator string, value interface{}, boolean ...string) QueryBuilder

	// WhereRaw 在查询中添加一个原生SQL“where”条件。
	//
	// sql: 原生SQL条件字符串。
	// bindings: SQL绑定参数数组。
	WhereRaw(sql string, bindings ...interface{}) QueryBuilder

	// OrWhereRaw 在查询中添加一个原生SQL“或where”条件。
	//
	// sql: 原生SQL条件字符串。
	// bindings: SQL绑定参数数组。
	OrWhereRaw(sql string, bindings ...interface{}) QueryBuilder

	// WhereIn 在指定列的值存在于给定的集合内时添加一个"where"条件。
	//
	// column: 要检查的列名。
	// values: 集合值。
	// not: 是否取反，默认为 false。
	WhereIn(column string, values []interface{}, not ...bool) QueryBuilder

	// Having 添加HAVING条件到查询语句中。
	//
	// column: 要应用HAVING条件的列名或原生SQL片段。
	// operator: 操作符，默认为 '='。
	// value: 条件值。
	// boolean: 条件之间的逻辑关系，默认为 'and'。
	Having(column interface{}, operator string, value interface{}, boolean ...string) QueryBuilder

	// OrHaving 添加一个“或HAVING”条件到查询语句中。
	//
	// column: 要应用HAVING条件的列名或原生SQL片段。
	// operator: 操作符，默认为 '='。
	// value: 条件值。
	OrHaving(column interface{}, operator string, value interface{}) QueryBuilder

	// Count 执行COUNT(*)聚合函数。
	//
	// column: （可选）要计数的特定列名，如果不提供则计算所有行数。
	Count(column ...string) (int64, error)

	// Max 执行MAX()聚合函数。
	//
	// column: 要查找最大值的列名。
	Max(column string) (float64, error)

	// Min 执行MIN()聚合函数。
	//
	// column: 要查找最小值的列名。
	Min(column string) (float64, error)

	// Avg 执行AVG()聚合函数。
	//
	// column: 要求平均值的列名。
	Avg(column string) (float64, error)

	// Union 将当前查询与另一个查询进行UNION操作。
	//
	// query: 另一个QueryBuilder实例。
	Union(query QueryBuilder) QueryBuilder

	// UnionAll 将当前查询与另一个查询进行UNION ALL操作。
	//
	// query: 另一个QueryBuilder实例。
	UnionAll(query QueryBuilder) QueryBuilder

	// Join 对数据库表进行JOIN操作。
	//
	// table: 要连接的表名。
	// first: 第一个连接条件（列名或原生SQL片段）。
	// operator: 连接条件的操作符，默认为 '='。
	// second: 第二个连接条件（列名或值）。
	// type: JOIN类型（如：INNER, LEFT, RIGHT等）。
	// where: （可选）附加的ON条件。
	Join(table string, first interface{}, operator string, second interface{}, types string, where ...interface{}) QueryBuilder

	// Distinct 在查询中添加 DISTINCT 关键字，以返回唯一结果。
	Distinct() QueryBuilder

	// Offset 设置查询的偏移量，用于分页操作。
	//
	// value: 要跳过的行数。
	Offset(value int) QueryBuilder

	// Limit 限制查询返回的结果数量。
	//
	// value: 要返回的最大记录数。
	Limit(value int) QueryBuilder

	// InsertGetId 执行插入操作并获取自增ID（如果存在）。
	//
	// values: 要插入的数据。
	// sequence: （可选）自增序列名称，默认根据数据库自动识别。
	InsertGetId(values interface{}, sequence ...string) (int64, error)

	// WhereExists 使用WHERE EXISTS子查询条件。
	//
	// closure: 返回QueryBuilder实例的闭包函数。
	WhereExists(closure func(QueryBuilder) QueryBuilder) QueryBuilder

	// WhereNotExists 使用WHERE NOT EXISTS子查询条件。
	//
	// closure: 返回QueryBuilder实例的闭包函数。
	WhereNotExists(closure func(QueryBuilder) QueryBuilder) QueryBuilder

	// FromSub 查询子查询作为FROM部分。
	//
	// query: 子查询生成器实例。
	// as: 给子查询起的别名。
	FromSub(query QueryBuilder, as string) QueryBuilder

	// JoinSub 进行JOIN操作时使用子查询。
	//
	// query: 子查询生成器实例。
	// as: 给子查询起的别名。
	// first: 第一个连接条件。
	// operator: 连接条件的操作符。
	// second: 第二个连接条件。
	// type: JOIN类型（如：INNER, LEFT, RIGHT等）。
	// where: （可选）附加的ON条件。
	JoinSub(query QueryBuilder, as string, first interface{}, operator string, second interface{}, types string, where ...interface{}) QueryBuilder

	// HavingRaw 添加原生SQL片段到HAVING子句中。
	//
	// sql: 原生SQL条件字符串。
	// bindings: SQL绑定参数数组。
	HavingRaw(sql string, bindings ...interface{}) QueryBuilder

	// OrHavingRaw 添加原生SQL片段到OR HAVING子句中。
	//
	// sql: 原生SQL条件字符串。
	// bindings: SQL绑定参数数组。
	OrHavingRaw(sql string, bindings ...interface{}) QueryBuilder

	// ToSql 获取当前查询构造器生成的SQL语句，不执行查询。
	//
	// bindings: 输出参数，将包含SQL绑定值。
	ToSql() (string, []interface{})

	// Exists 检查是否有与给定查询相匹配的记录存在。
	Exists() (bool, error)

	// GetBindings 获取查询构建器中的所有绑定值。
	GetBindings() []interface{}

	// Chunk 分块处理查询结果集。
	//
	// count: 每次处理的记录数。
	// callback: 对每批次数据执行的回调函数。
	Chunk(count int, callback func([]Record) bool) error

	// Tap 允许在不改变链式调用的情况下对查询对象进行中间操作。
	//
	// callback: 对当前QueryBuilder实例执行的操作函数。
	Tap(callback func(QueryBuilder)) QueryBuilder

	// WhereNull 指定列的值为 NULL 时添加一个"where"条件。
	//
	// column: 列名。
	WhereNull(column string) QueryBuilder

	// WhereNotNull 指定列的值不为 NULL 时添加一个"where"条件。
	//
	// column: 列名。
	WhereNotNull(column string) QueryBuilder

	// OrWhereLike 添加一个“或where LIKE”条件到查询语句中。
	//
	// column: 要进行模糊匹配的列名。
	// value: 包含通配符（%）的匹配字符串。
	OrWhereLike(column string, value interface{}) QueryBuilder

	// WhereInColumns 在指定列的值存在于另一列的集合内时添加一个"where"条件。
	//
	// first: 第一个列名。
	// operator: 连接操作符，默认是 'in' 或 'not in'。
	// second: 第二个列名，或者是一个包含列值的数组。
	WhereInColumns(first string, operator string, second interface{}) QueryBuilder

	// OrWhereBetween 添加一个“或where BETWEEN”条件到查询语句中。
	//
	// column: 要检查范围的列名。
	// values: 包含范围起始和结束值的数组。
	OrWhereBetween(column string, values []interface{}) QueryBuilder

	// WhereColumn 在一个列的值与另一个列的值比较时添加一个"where"条件。
	//
	// first: 第一个列名或原生SQL片段。
	// operator: 比较运算符。
	// second: 第二个列名或原生SQL片段。
	WhereColumn(first interface{}, operator string, second interface{}) QueryBuilder

	// OrWhereColumn 在一个列的值与另一个列的值比较时添加一个“或where”条件。
	//
	// first: 第一个列名或原生SQL片段。
	// operator: 比较运算符。
	// second: 第二个列名或原生SQL片段。
	OrWhereColumn(first interface{}, operator string, second interface{}) QueryBuilder

	// When 对给定的布尔条件执行相应的查询构建操作。
	//
	// condition: 返回布尔值的闭包函数。
	// callback: 当条件满足时执行的查询构建回调函数。
	When(condition func() bool, callback func(QueryBuilder) QueryBuilder) QueryBuilder

	// WhereLike 在指定列进行模糊匹配时添加一个"where"条件。
	//
	// column: 要进行模糊匹配的列名。
	// value: 包含通配符（%）的匹配字符串。
	WhereLike(column string, value interface{}) QueryBuilder

	// OrWhereNotNull 指定列的值不为 NULL 时添加一个“或where”条件。
	//
	// column: 列名。
	OrWhereNotNull(column string) QueryBuilder

	// WhereNotIn 在指定列的值不在给定集合内时添加一个"where"条件。
	//
	// column: 要检查的列名。
	// values: 集合值数组。
	WhereNotIn(column string, values []interface{}) QueryBuilder

	// OrWhereNotIn 在指定列的值不在给定集合内时添加一个“或where”条件。
	//
	// column: 要检查的列名。
	// values: 集合值数组。
	OrWhereNotIn(column string, values []interface{}) QueryBuilder

	// WhereDate 对日期部分进行比较时添加一个"where"条件。
	//
	// column: 要比较的日期列名。
	// operator: 比较运算符。
	// value: 日期值。
	WhereDate(column string, operator string, value interface{}) QueryBuilder

	// OrWhereDate 对日期部分进行比较时添加一个“或where”条件。
	//
	// column: 要比较的日期列名。
	// operator: 比较运算符。
	// value: 日期值。
	OrWhereDate(column string, operator string, value interface{}) QueryBuilder

	// WhereTime 对时间部分进行比较时添加一个"where"条件。
	//
	// column: 要比较的时间列名。
	// operator: 比较运算符。
	// value: 时间值。
	WhereTime(column string, operator string, value interface{}) QueryBuilder

	// OrWhereTime 对时间部分进行比较时添加一个“或where”条件。
	//
	// column: 要比较的时间列名。
	// operator: 比较运算符。
	// value: 时间值。
	OrWhereTime(column string, operator string, value interface{}) QueryBuilder

	// WhereMonth 对月份部分进行比较时添加一个"where"条件。
	//
	// column: 要比较的日期列名。
	// operator: 比较运算符。
	// value: 月份值。
	WhereMonth(column string, operator string, value interface{}) QueryBuilder

	// OrWhereMonth 对月份部分进行比较时添加一个“或where”条件。
	//
	// column: 要比较的日期列名。
	// operator: 比较运算符。
	// value: 月份值。
	OrWhereMonth(column string, operator string, value interface{}) QueryBuilder

	// WhereYear 对年份部分进行比较时添加一个"where"条件。
	//
	// column: 要比较的日期列名。
	// operator: 比较运算符。
	// value: 年份值。
	WhereYear(column string, operator string, value interface{}) QueryBuilder

	// OrWhereYear 对年份部分进行比较时添加一个“或where”条件。
	//
	// column: 要比较的日期列名。
	// operator: 比较运算符。
	// value: 年份值。
	OrWhereYear(column string, operator string, value interface{}) QueryBuilder

	// WhereDay 对日期中的日部分进行比较时添加一个"where"条件。
	//
	// column: 要比较的日期列名。
	// operator: 比较运算符。
	// value: 日值。
	WhereDay(column string, operator string, value interface{}) QueryBuilder

	// OrWhereDay 对日期中的日部分进行比较时添加一个“或where”条件。
	//
	// column: 要比较的日期列名。
	// operator: 比较运算符。
	// value: 日值。
	OrWhereDay(column string, operator string, value interface{}) QueryBuilder

	// WhereColumnEqual 当两个列的值相等时添加一个"where"条件。
	//
	// first: 第一个列名。
	// second: 第二个列名。
	WhereColumnEqual(first string, second string) QueryBuilder

	// OrWhereColumnEqual 当两个列的值相等时添加一个“或where”条件。
	//
	// first: 第一个列名。
	// second: 第二个列名。
	OrWhereColumnEqual(first string, second string) QueryBuilder

	// WhereJsonContains 在JSON列中查找给定的值时添加一个"where"条件。
	//
	// column: JSON列名。
	// value: 要在JSON列中查找的值。
	// path: （可选）指定要搜索的JSON路径。
	WhereJsonContains(column string, value interface{}, path ...string) QueryBuilder

	// OrWhereJsonContains 在JSON列中查找给定的值时添加一个“或where”条件。
	//
	// column: JSON列名。
	// value: 要在JSON列中查找的值。
	// path: （可选）指定要搜索的JSON路径。
	OrWhereJsonContains(column string, value interface{}, path ...string) QueryBuilder

	// WhereJsonLength 对JSON列中的数组或对象长度进行比较时添加一个"where"条件。
	//
	// column: JSON列名。
	// operator: 比较运算符（如：'>', '<=', '=' 等）。
	// value: 要比较的长度值。
	// path: （可选）指定要检查长度的JSON路径。
	WhereJsonLength(column string, operator string, value interface{}, path ...string) QueryBuilder

	// OrWhereJsonLength 对JSON列中的数组或对象长度进行比较时添加一个“或where”条件。
	//
	// column: JSON列名。
	// operator: 比较运算符（如：'>', '<=', '=' 等）。
	// value: 要比较的长度值。
	// path: （可选）指定要检查长度的JSON路径。
	OrWhereJsonLength(column string, operator string, value interface{}, path ...string) QueryBuilder

	// OrWhereExists 使用WHERE EXISTS子查询条件。
	//
	// callback: 返回QueryBuilder实例的闭包函数，用于构建子查询。
	OrWhereExists(callback func(QueryBuilder) QueryBuilder) QueryBuilder

	// OrWhereNotExists 使用WHERE NOT EXISTS子查询条件。
	//
	// callback: 返回QueryBuilder实例的闭包函数，用于构建子查询。
	OrWhereNotExists(callback func(QueryBuilder) QueryBuilder) QueryBuilder

	// CrossJoin 进行CROSS JOIN操作。
	//
	// table: 要连接的表名。
	// first: （可选）连接条件。
	// operator: （可选）连接条件的操作符。
	// second: （可选）连接条件的值。
	CrossJoin(table string, first interface{}, operator string, second interface{}) QueryBuilder

	// LeftJoinSub 使用子查询进行LEFT JOIN操作，并给子查询指定别名。
	//
	// subquery: 子查询生成器实例。
	// as: 给子查询起的别名。
	// first: 第一个连接条件（列名或表达式）。
	// operator: 连接条件的操作符，如 '='。
	// second: 第二个连接条件（列名或值）。
	// where: （可选）附加的ON条件。
	LeftJoinSub(subquery QueryBuilder, as string, first interface{}, operator string, second interface{}, where ...interface{}) QueryBuilder

	// RightJoinSub 使用子查询进行RIGHT JOIN操作，并给子查询指定别名。
	//
	// subquery: 子查询生成器实例。
	// as: 给子查询起的别名。
	// first: 第一个连接条件（列名或表达式）。
	// operator: 连接条件的操作符，如 '='。
	// second: 第二个连接条件（列名或值）。
	// where: （可选）附加的ON条件。
	RightJoinSub(subquery QueryBuilder, as string, first interface{}, operator string, second interface{}, where ...interface{}) QueryBuilder

	// FullJoinSub 使用子查询进行FULL OUTER JOIN操作，并给子查询指定别名。
	//
	// subquery: 子查询生成器实例。
	// as: 给子查询起的别名。
	// first: 第一个连接条件（列名或表达式）。
	// operator: 连接条件的操作符，如 '='。
	// second: 第二个连接条件（列名或值）。
	// where: （可选）附加的ON条件。
	FullJoinSub(subquery QueryBuilder, as string, first interface{}, operator string, second interface{}, where ...interface{}) QueryBuilder

	// Intersect 通过INTERSECT操作将当前查询与其他查询的结果集相交。
	//
	// query: 另一个QueryBuilder实例。
	Intersect(query QueryBuilder) QueryBuilder

	// Except 通过EXCEPT操作从当前查询中排除另一个查询的结果集。
	//
	// query: 另一个QueryBuilder实例。
	Except(query QueryBuilder) QueryBuilder

	// From 添加FROM子句，指定表名或子查询。
	From(table string) QueryBuilder

	// Joins 添加多个JOIN条件到查询语句中。
	//
	// joins: JOIN条件的数组，每个元素是一个JoinClause结构或其他表示JOIN条件的数据结构。
	Joins(joins ...interface{}) QueryBuilder

	// JoinClause 创建一个JOIN子句实例，用于更复杂的JOIN操作构造。
	JoinClause(table string, first interface{}, operator string, second interface{}, types string) JoinClause

	// LeftJoin 执行LEFT JOIN操作。
	LeftJoin(table string, first interface{}, operator string, second interface{}, where ...interface{}) QueryBuilder

	// RightJoin 执行RIGHT JOIN操作。
	RightJoin(table string, first interface{}, operator string, second interface{}, where ...interface{}) QueryBuilder

	// DistinctColumn 在SELECT语句中使用DISTINCT关键字对指定列进行去重。
	DistinctColumn(columns ...string) QueryBuilder

	// FirstOrFail 获取第一条记录，如果没有记录则抛出错误。
	FirstOrFail() (Record, error)

	// Pluck 从查询结果集中获取指定列的值列表。
	Pluck(column string) ([]interface{}, error)

	// ToSqlWithBindings 返回生成的SQL语句以及绑定参数。
	ToSqlWithBindings() (string, []interface{})

	// List 获取指定列的键值对列表。
	List(column string, keyColumn string) ([]map[string]interface{}, error)

	// WithEagerLoads 添加预加载关联模型到查询中。
	WithEagerLoads(relations ...string) QueryBuilder

	// InsertIgnore 类似Insert，但在某些数据库（如MySQL）中支持INSERT IGNORE语句。
	InsertIgnore(values interface{}) (int64, error)

	// UpdateOrInsert 如果记录存在则更新，不存在则插入数据。
	UpdateOrInsert(attributes map[string]interface{}, values interface{}) (bool, int64, error)

	// ToRawSql 将查询转换为原生SQL字符串，并用？占位符替换所有绑定值。
	ToRawSql() string

	// WhereInJson 当JSON列的某个路径下的值存在于给定集合内时添加一个"where"条件。
	WhereInJson(column string, values interface{}, path ...string) QueryBuilder

	// OrWhereInJson 当JSON列的某个路径下的值存在于给定集合内时添加一个“或where”条件。
	OrWhereInJson(column string, values interface{}, path ...string) QueryBuilder

	// WhereNotInJson 当JSON列的某个路径下的值不在给定集合内时添加一个"where"条件。
	WhereNotInJson(column string, values interface{}, path ...string) QueryBuilder

	// OrWhereNotInJson 当JSON列的某个路径下的值不在给定集合内时添加一个“或where”条件。
	OrWhereNotInJson(column string, values interface{}, path ...string) QueryBuilder

	// WhereJsonContainsPath 当JSON列包含指定路径及对应的值时添加一个"where"条件。
	WhereJsonContainsPath(column string, value interface{}, path ...string) QueryBuilder

	// OrWhereJsonContainsPath 当JSON列包含指定路径及对应的值时添加一个“或where”条件。
	OrWhereJsonContainsPath(column string, value interface{}, path ...string) QueryBuilder

	// ForPage 实现分页查询，基于页码和每页数量。
	ForPage(page int, perPage int) QueryBuilder

	// Latest 或 Oldest 根据指定列按降序或升序获取最新或最旧的一条记录。
	Latest(column string) QueryBuilder
	Oldest(column string) QueryBuilder

	// WithLock 加锁查询以防止并发修改。
	WithLock(lockType string) QueryBuilder // lockType 可能是 "for update", "lock in share mode" 等根据数据库方言确定的锁定模式。

	// ExistsSub 查询子查询作为EXISTS条件。
	ExistsSub(query QueryBuilder) QueryBuilder

	// Raw 对查询中的一部分或全部使用原始SQL表达式。
	Raw(sql string, bindings ...interface{}) QueryBuilder
}

// JoinClause 是一个辅助结构体，用于构建复杂的JOIN条件。
type JoinClause struct {
	Table    string
	First    interface{}
	Operator string

	Second     interface{}
	Type       string
	On         interface{}
	Wheres     []interface{}
	Columns    []string
	Alias      string
	InnerJoins []JoinClause // 其他JoinClause内部实现细节...
}

// ExampleUsage:
// qb := NewQueryBuilder()
// records, err := qb.Table("users").Select("name", "email").Where("age", ">=", 18).OrderBy("name", "asc").Get()
