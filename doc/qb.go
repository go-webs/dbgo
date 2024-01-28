package doc

import "database/sql"

// 其他可能的方法例如：
func (qb *QueryBuilder) Table(tableName string) *QueryBuilder
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder
func (qb *QueryBuilder) Where(cond Column, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) Join(joinType string, table string, firstCol string, operator string, secondCol string, on ...interface{}) *QueryBuilder
func (qb *QueryBuilder) OrderBy(column string, direction string) *QueryBuilder
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder

// QueryBuilder 方法继续补充：
func (qb *QueryBuilder) WhereRaw(sql string, bindings ...interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereRaw(sql string, bindings ...interface{}) *QueryBuilder
func (qb *QueryBuilder) HavingRaw(sql string, bindings ...interface{}) *QueryBuilder
func (qb *QueryBuilder) OrHavingRaw(sql string, bindings ...interface{}) *QueryBuilder
func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder
func (qb *QueryBuilder) Having(conditions ...Condition) *QueryBuilder
func (qb *QueryBuilder) OrHaving(conditions ...Condition) *QueryBuilder
func (qb *QueryBuilder) UnionAll(query *QueryBuilder) *QueryBuilder
func (qb *QueryBuilder) Union(query *QueryBuilder) *QueryBuilder
func (qb *QueryBuilder) Insert(data map[string]interface{}) (int64, error)
func (qb *QueryBuilder) Update(data map[string]interface{}, where Conditions) (int64, error)
func (qb *QueryBuilder) Delete(where Conditions) (int64, error)
func (qb *QueryBuilder) Pluck(column string) ([]interface{}, error)
func (qb *QueryBuilder) Paginate(page int, perPage int) (*Paginator, error)

// Exists 查询记录是否存在并返回布尔值。
func (qb *QueryBuilder) Exists() (bool, error)

// Count 执行COUNT(*)聚合操作并返回计数。
func (qb *QueryBuilder) Count() (int64, error)

// Distinct 添加DISTINCT关键字到SELECT子句中。
func (qb *QueryBuilder) Distinct(columns ...string) *QueryBuilder

// SelectSub 选择一个子查询作为列。
func (qb *QueryBuilder) SelectSub(subQuery *QueryBuilder, alias string) *QueryBuilder

// JoinSub 使用子查询进行JOIN操作。
func (qb *QueryBuilder) JoinSub(subQuery *QueryBuilder, asAlias string, onClause Column, operator string, refColumn string) *QueryBuilder

// Raw 在查询语句中插入原生SQL片段。
func (qb *QueryBuilder) Raw(sqlFragment string, bindings ...interface{}) *QueryBuilder

// ForUpdate 设置查询以获取行级锁定（FOR UPDATE）。
func (qb *QueryBuilder) ForUpdate() *QueryBuilder

// LockShared 设置查询以获取共享锁定（LOCK IN SHARE MODE）。
func (qb *QueryBuilder) LockShared() *QueryBuilder

// WithTrashed 包含已删除的数据（适用于软删除场景）。
func (qb *QueryBuilder) WithTrashed() *QueryBuilder

// WithoutTrashed 不包含已删除的数据（适用于软删除场景）。
func (qb *QueryBuilder) WithoutTrashed() *QueryBuilder

// QueryBuilder 方法补充：
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereIn(column string, values []interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereNotIn(column string, values []interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereNotIn(column string, values []interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereBetween(column string, from interface{}, to interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereBetween(column string, from interface{}, to interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereNull(columns ...string) *QueryBuilder
func (qb *QueryBuilder) OrWhereNull(columns ...string) *QueryBuilder
func (qb *QueryBuilder) WhereNotNull(columns ...string) *QueryBuilder
func (qb *QueryBuilder) OrWhereNotNull(columns ...string) *QueryBuilder
func (qb *QueryBuilder) WhereColumn(column string, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereColumn(column string, operator string, value interface{}) *QueryBuilder

// JSON相关操作：
func (qb *QueryBuilder) WhereJsonContains(column string, value interface{}, path ...string) *QueryBuilder
func (qb *QueryBuilder) OrWhereJsonContains(column string, value interface{}, path ...string) *QueryBuilder
func (qb *QueryBuilder) WhereJsonLength(column string, operator string, value interface{}, path ...string) *QueryBuilder
func (qb *QueryBuilder) OrWhereJsonLength(column string, operator string, value interface{}, path ...string) *QueryBuilder
func (qb *QueryBuilder) WhereJsonPath(column string, operator string, value interface{}, path ...string) *QueryBuilder
func (qb *QueryBuilder) OrWhereJsonPath(column string, operator string, value interface{}, path ...string) *QueryBuilder

// 时间日期比较：
func (qb *QueryBuilder) WhereDate(column string, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereDate(column string, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereTime(column string, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereTime(column string, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereMonth(column string, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereMonth(column string, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereYear(column string, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereYear(column string, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereDay(column string, operator string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereDay(column string, operator string, value interface{}) *QueryBuilder

// 聚合函数：
func (qb *QueryBuilder) Sum(column string) (float64, error)
func (qb *QueryBuilder) Max(column string) (interface{}, error)
func (qb *QueryBuilder) Min(column string) (interface{}, error)
func (qb *QueryBuilder) Avg(column string) (float64, error)

// 更多高级查询：
func (qb *QueryBuilder) ExistsSub(query func(*QueryBuilder) *QueryBuilder) *QueryBuilder
func (qb *QueryBuilder) WhereExists(query func(*QueryBuilder) *QueryBuilder) *QueryBuilder
func (qb *QueryBuilder) WhereNotExists(query func(*QueryBuilder) *QueryBuilder) *QueryBuilder

// 其他实用方法：
func (qb *QueryBuilder) ToSQL() (string, []interface{})
func (qb *QueryBuilder) Clone() *QueryBuilder
func (qb *QueryBuilder) ResetQueryParts() *QueryBuilder

// JoinClause 结构体补充方法：
func (jc *JoinClause) OnRaw(sql string, bindings ...interface{}) *JoinClause
func (jc *JoinClause) OrOnRaw(sql string, bindings ...interface{}) *JoinClause

// 创建和执行子查询
func (qb *QueryBuilder) SubQuery() *SubQueryBuilder

// QueryBuilder 方法继续补充：
func (qb *QueryBuilder) DistinctColumn(columns ...string) *QueryBuilder
func (qb *QueryBuilder) HavingBetween(column string, from interface{}, to interface{}) *QueryBuilder
func (qb *QueryBuilder) OrHavingBetween(column string, from interface{}, to interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereNotLike(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereNotLike(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereILike(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereILike(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereNotILike(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) OrWhereNotILike(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WithRelations(relations []string) *QueryBuilder                                   // 预加载关联数据（适用于ORM库）
func (qb *QueryBuilder) ForPage(page int, perPage int) (*Paginator, error)                                // 分页查询并返回分页器对象
func (qb *QueryBuilder) ToStatement() (*sql.Stmt, error)                                                  // 将构建的查询语句预编译为SQL Statement
func (qb *QueryBuilder) ToNamedQuery(name string) (*NamedQuery, error)                                    // 将查询保存为可重用的命名查询
func (qb *QueryBuilder) InsertGetID(values map[string]interface{}, sequenceName ...string) (int64, error) // 插入数据并获取自增ID
func (qb *QueryBuilder) InsertIgnore(values map[string]interface{}) (int64, error)                        // 执行INSERT IGNORE操作
func (qb *QueryBuilder) Truncate() error                                                                  // 清空表数据

func (sb *SubQueryBuilder) Select(columns ...string) *SubQueryBuilder
func (sb *SubQueryBuilder) Where(...interface{}) *SubQueryBuilder
func (sb *SubQueryBuilder) ToQueryBuilder() *QueryBuilder

// QueryBuilder 方法继续补充：
func (qb *QueryBuilder) HavingIn(column string, values []interface{}) *QueryBuilder
func (qb *QueryBuilder) OrHavingIn(column string, values []interface{}) *QueryBuilder
func (qb *QueryBuilder) HavingNotIn(column string, values []interface{}) *QueryBuilder
func (qb *QueryBuilder) OrHavingNotIn(column string, values []interface{}) *QueryBuilder

// 设置表别名
func (qb *QueryBuilder) TableAlias(alias string) *QueryBuilder

// 查询构建器中使用CASE表达式
func (qb *QueryBuilder) Case(whenClauses ...*WhenClause) *CaseStatement

// 执行CASE表达式并将其添加到查询中
func (cs *CaseStatement) End() *QueryBuilder

// 高级查询构造辅助方法
func (qb *QueryBuilder) WhereInSub(query func(*QueryBuilder) *QueryBuilder, column string) *QueryBuilder
func (qb *QueryBuilder) WhereNotInSub(query func(*QueryBuilder) *QueryBuilder, column string) *QueryBuilder

// 聚合函数扩展（根据数据库驱动支持）
func (qb *QueryBuilder) BitAnd(column string) *QueryBuilder
func (qb *QueryBuilder) BitOr(column string) *QueryBuilder
func (qb *QueryBuilder) BitXor(column string) *QueryBuilder
func (qb *QueryBuilder) StdDev(column string) (float64, error)
func (qb *QueryBuilder) Variance(column string) (float64, error)

// 针对时间字段的查询方法
func (qb *QueryBuilder) WhereYearEquals(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereMonthEquals(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereDayOfMonthEquals(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereHourEquals(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereMinuteEquals(column string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereSecondEquals(column string, value interface{}) *QueryBuilder

// 锁定行或表
func (qb *QueryBuilder) ForShare() *QueryBuilder
func (qb *QueryBuilder) LockForUpdate() *QueryBuilder

// 其他高级操作（如窗口函数、CTE等，具体实现取决于数据库驱动的支持）
func (qb *QueryBuilder) With(name string, subquery func(*QueryBuilder) *QueryBuilder) *QueryBuilder
func (qb *QueryBuilder) Over(partitionBy ...string) *WindowFunctionBuilder

// QueryBuilder 方法继续补充：

// 排序支持多个字段，每个字段可以有自己的排序方向
func (qb *QueryBuilder) OrderByMulti(columns []OrderByItem) *QueryBuilder

// 复制查询构建器以创建新的查询实例
func (qb *QueryBuilder) Clone() *QueryBuilder

// 根据列值进行升序或降序分组
func (qb *QueryBuilder) GroupByAsc(column string) *QueryBuilder
func (qb *QueryBuilder) GroupByDesc(column string) *QueryBuilder

// 创建窗口函数并添加到查询中
func (wfb *WindowFunctionBuilder) Build() (*QueryBuilder, error)
func (wfb *WindowFunctionBuilder) RowsBetween(start interface{}, end interface{}) *WindowFunctionBuilder
func (wfb *WindowFunctionBuilder) RangeBetween(start interface{}, end interface{}) *WindowFunctionBuilder
func (wfb *WindowFunctionBuilder) FrameType(frameType string) *WindowFunctionBuilder

// 链式调用设置表别名
func (qb *QueryBuilder) As(alias string) *QueryBuilder

// 查询构造器中的嵌套集合操作（如Oracle的CONNECT BY和MySQL的WITH ... AS）
func (qb *QueryBuilder) WithRecursive(name string, columns []string, baseQuery func(*QueryBuilder) *QueryBuilder, recursiveQuery func(*QueryBuilder) *QueryBuilder) *QueryBuilder

// 支持JSON操作
func (qb *QueryBuilder) JsonGet(jsonColumn string, path string) *QueryBuilder
func (qb *QueryBuilder) JsonSet(jsonColumn string, path string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) JsonInsert(jsonColumn string, path string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) JsonReplace(jsonColumn string, path string, value interface{}) *QueryBuilder
func (qb *QueryBuilder) JsonRemove(jsonColumn string, path string) *QueryBuilder

// 使用表达式作为查询条件或排序依据
func (qb *QueryBuilder) WhereRawExp(expression Expression) *QueryBuilder
func (qb *QueryBuilder) OrderByRawExp(expression Expression) *QueryBuilder

// 其他数据库特定功能，例如PostgreSQL中的数组操作
func (qb *QueryBuilder) ArrayContains(column string, values []interface{}) *QueryBuilder
func (qb *QueryBuilder) ArrayOverlaps(column string, values []interface{}) *QueryBuilder

// QueryBuilder 方法继续补充：

// 分页查询并返回分页结果集（包含数据和分页信息）
func (qb *QueryBuilder) PaginateWithInfo(page int, perPage int) (*PagingResult, error)

// 高级子查询支持，例如用于构建自关联查询或嵌套查询
func (qb *QueryBuilder) SubQueryAs(alias string, subqueryFunc func(subQb *QueryBuilder) *QueryBuilder) *QueryBuilder

// 批量插入数据
func (qb *QueryBuilder) InsertBatch(records []map[string]interface{}) (int64, error)

// 更新查询中的特定字段，仅更新非空值
func (qb *QueryBuilder) UpdateNonEmpty(values map[string]interface{}) (int64, error)

// 查询构造器中处理字符串连接函数
func (qb *QueryBuilder) Concatenate(columns ...string) *QueryBuilder

// 使用表达式构建查询条件时支持参数绑定
func (qb *QueryBuilder) WhereExp(expression Expression, operator string, value interface{}) *QueryBuilder

// 支持全文搜索功能（如MySQL的MATCH AGAINST、PostgreSQL的tsvector操作等）
func (qb *QueryBuilder) FullTextSearch(column string, terms string, booleanMode bool) *QueryBuilder

// 处理多表删除（Cascading Delete）
func (qb *QueryBuilder) DeleteRelated(table string, foreignKey string, onDeleteAction string) *QueryBuilder

// 其他数据库特性支持，比如PostgreSQL的JSONB操作
func (qb *QueryBuilder) JsonbExtract(column string, path string) *QueryBuilder
func (qb *QueryBuilder) JsonbExists(column string, path string) *QueryBuilder

// QueryBuilder 方法继续补充：

// 执行原生SQL查询并返回结果集
func (qb *QueryBuilder) RawQuery(sql string, bindings ...interface{}) (*sql.Rows, error)

// 使用事务执行查询或操作
func (qb *QueryBuilder) Transactional(transactionFunc func(*QueryBuilder) error) error

// 创建一个临时表并插入数据，用于后续查询
func (qb *QueryBuilder) CreateTemporaryTable(tableName string, columns []string, data [][]interface{}) error

// 查询构造器中的窗口函数支持
func (qb *QueryBuilder) Over(partitionBy []string, orderBy []OrderByItem, windowFrame WindowFrame) *WindowFunctionQueryBuilder

type WindowFunctionQueryBuilder struct {
	BaseQueryBuilder *QueryBuilder
	PartitionBy      []string
	OrderBy          []OrderByItem
	Frame            WindowFrame
}

func (wfqb *WindowFunctionQueryBuilder) Function(functionName string) *WindowFunctionQueryBuilder
func (wfqb *WindowFunctionQueryBuilder) Build() (*QueryBuilder, error)

// 对于支持的数据库，实现窗口函数排序
func (qb *QueryBuilder) RowNumberOver(orderBy []OrderByItem) *QueryBuilder
func (qb *QueryBuilder) RankOver(orderBy []OrderByItem) *QueryBuilder
func (qb *QueryBuilder) DenseRankOver(orderBy []OrderByItem) *QueryBuilder

// 支持数据库的特定函数调用
func (qb *QueryBuilder) CallDatabaseFunction(functionName string, args ...interface{}) *QueryBuilder

// 处理时间序列分析相关的函数（如LAG、LEAD、FIRST_VALUE、LAST_VALUE等）
func (qb *QueryBuilder) Lag(column string, offset int, defaultVal interface{}) *QueryBuilder
func (qb *QueryBuilder) Lead(column string, offset int, defaultVal interface{}) *QueryBuilder
func (qb *QueryBuilder) FirstValue(column string) *QueryBuilder
func (qb *QueryBuilder) LastValue(column string) *QueryBuilder

// 针对特定数据库的JOIN类型扩展（如Oracle的OUTER APPLY、CROSS APPLY）
func (qb *QueryBuilder) Apply(table string, alias string, condition *WhereClause) *QueryBuilder
