# dbgo
Php Laravel orm Eloquent 的 go 实现, 与官方文档保持一致
https://laravel.com/docs/10.x/queries

## 建表
```go
package main

import (
	"time"
	// 一定要引入数据库驱动, 这是dbgo提供的mysql驱动,包括orm解析都在这里
	// 其他数据库驱动, 可以自行实现, 按照 dbog.IDriver 接口实现接口
	_ "github.com/go-webs/dbgo-driver-mysql"
)

type Users struct {
	Id        int       `db:"id,pk"`    // 设定字段名字为 id, pk意为该字段为主键 primary key
	Name      string    `db:"name"`     // 设定字段名字为 name, 非主键不需要做任何标记
	Email     string    `db:"email"`
	Title     string    `db:"title"`
	Active    bool      `db:"active"`
	Votes     int       `db:"votes"`
	Balance   float64   `db:"balance"`
	CreatedAt time.Time `db:"created_at"`   // datetime 类型, 记得要在连接的dsn后边加上 parseTime=true, 见下边 dbgo.Open() 示例
}
// TableName 手动指定 Users struct 表名为 users, 如果不指定, 则自动解析为 struct 名字, 即(Users)
func (Users) TableName() string {
    return "users"
}
var dg = dbgo.Open("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true")
// 实际项目中,自定义helper函数,方便调用orm,包外调用,可以设置成大写导出的函数使用
func db() dbgo.Database {
	return dg.NewDatabase()
}
```

## 原生sql
```go
s, err := db().NewSession()
s.Query("select * from users limit ?", 10)
s.Exec("delete from users where id=?", 1)

s.Transaction(func(db dbgo.Session) error {
	db.Exec("insert into users (name, email) values (?, ?)", "张三", "aa@aa.com")
    // 这里可以写多个sql语句, 事务执行
	...
	// error case
	db.Rollback()
	...
	// finish
	db.Commit()
}
```

## go style
```go
// select xx,xx from users limit 10
var users []Users
db().Limit(10).To(&users)

// select xx,xx from users limit 1
var user Users
db().To(&users)

// select xx,xx from users where id=1
var user = User{Id: 1}
db().To(&user)
// delete from users where id=1
db().Delete(&user)

// insert into users (name, email) values ("张三", "aa@aa.com")
var user = User{Name: "张三", Email: "aa@aa.com"}
db().Insert(&user)

// update users set name="李四" where id=1
var user = User{Id: 1, Name:"李四"}
db().Update(&user)

// 自动事务
db().Transaction(func(db dbgo.Database) error {
	db.Insert(&user)
    db.Update(&user)
	db.To(&user)
}
// 手动事务
var tx = db()
tx.Begin()
tx.Rollback()
tx.Commit()

// 自动嵌套事务
db().Transaction(func(db dbgo.Database) error {
    db().Transaction(func(db dbgo.Database) error {
    }
}
// 手动嵌套事务
var tx = db()

tx.Begin()

// 自动子事务
tx.Begin() // 自动 savepoint 子事务
tx.Rollback()   // 自动回滚到上一个 savepoint
// 手动子事务
tx.SavePoint("savepoint1")    // 手动 savepoint 到 savepoint1(自定义名字)
tx.RollbackTo("savepoint1") // 手动回滚到自定义的 savepoint

tx.Commit()
```
go style 可以使用下边 php style 的所有条件方法, 如: `join(), where(), having(), order(), limit(), offset()`

## php laravel style
### Running Database Queries
```go
// db().Table().Select().Where().GroupBY().Having().OrderBy().Limit().Offset()
// Retrieving All Rows From a Table

db().Table("users").Get()
// 等同于
db().Table(Users{}).Get()    // Users{} 等同于 "users", 都可以作为表名, orm会自动识别并解析出设定的表名

// Retrieving a Single Row / Column From a Table
db().Table("users").Where("name", "John").First()
db().Table("users").Value("email")
db().Table("users").Find(3)
// Retrieving a List of Column Values
db().Table("users").Pluck("title")
db().Table("users").List("title", "name")    // {name: title}
// Chunking Results
db().Table("users").OrderBy("id").Chunk(100, func([]Users){/* some codes */})
db().Table("users").Where("active", false).ChunkById(100, func([]Users){/* some codes */})
// Streaming Results Lazily
db().Table("users").OrderBy("id").Lazy().Each(func(Users){/* some codes */})
db().Table("users").Where("active", false).LazyById().Each(func(Users){/* some codes */})
// Aggregates
db().Table("users").Count()
db().Table("orders").Max("price")
db().Table("orders").Where("finalized", 1).Avg("price")
// Determining if Records Exist
db().Table("orders").Where("finalized", 1).Exists()
db().Table("users").Where("finalized", 1).DoesntExist()
```
### Select Statements
```go
db().Table("users").Select("name", "email as user_email").Get()
db().Table("users").Distinct().Get()
var query = db().Table("users").Select("name")
query.AddSelect("age").Get()
```
### Raw Expressions
```go
db().Table("users").
	Select(dbgo.Raw("count(*) as user_count, status")).
	Where("status", "<>", 1).
	GroupBy("status").
	Get()
```
#### Raw Methods
- selectRaw
```go
db().Table("orders").
	SelectRaw("price * ? as price_with_tax", [1.0825]).
	Get()
```
- whereRaw / orWhereRaw
```go
db().Table("orders").
	WhereRaw("price > IF(state = "TX", ?, 100)", [200]).
	Get()
```
- havingRaw / orHavingRaw
```go
db().Table("orders").
	Select("department", dbgo.Raw("SUM(price) as total_sales")).
	GroupBy("department").
	HavingRaw("SUM(price) > ?", [2500]).
	Get()
```
- orderByRaw
```go
db().Table("orders").
	OrderByRaw("updated_at - created_at DESC").
	Get()
```
- groupByRaw
```go
db().Table("orders").
	Select("city", "state").
	GroupByRaw("city, state").
	Get()
```
### Joins
- Inner Join Clause
```go
db().Table("users").
	Join("contacts", "users.id", "=", "contacts.user_id").
	Join("orders", "users.id", "=", "orders.user_id").
	Select("users.*", "contacts.phone", "orders.price").
	Get()
```
- Left Join / Right Join Clause
```go
db().Table("users").
    LeftJoin("posts", "users.id", "=", "posts.user_id").
    Get()

db().Table("users").
    RightJoin("posts", "users.id", "=", "posts.user_id").
    Get()
```
- Cross Join Clause
```go
db().Table("sizes").
    CrossJoin("colors").
    Get()
```
- Advanced Join Clauses
```go
db().Table("users").Join("contacts", func (joins dbgo.JoinClause) {
            joins.On("users.id", "=", "contacts.user_id").OrOn(/* ... */);
        }).
        Get()

//db().Table("users").
//    Join("contacts", func (joins dbgo.JoinClause) {
//        joins.On("users.id", "=", "contacts.user_id").Where("contacts.user_id", ">", 5)
//    }).
//    Get();
```
- Subquery Joins
```go
//latestPosts := db().Table("posts").Select("user_id", dbgo.Raw("MAX(created_at) as last_post_created_at")).Where("is_published", true).GroupBy("user_id")
//db().Table("users").JoinSub(latestPosts, "latest_posts", function (joins dbgo.JoinClause) {
//            joins.On("users.id", "=", "latest_posts.user_id")
//        }).Get()
```

### Unions
```go
first := db().Table("users").WhereNull("first_name")
//db().Table("users").WhereNull("last_name").Union(first).Get()
db().Table("users").WhereNull("last_name").Union(first)
```
### Basic Where Clauses
```go
// Where Clauses
db().Table("users").Where("votes","=",100).Where("age",">",35).Get()
db().Table("users").Where("votes",">=",100).Get()
db().Table("users").Where("votes","<>",100).Get()
db().Table("users").Where("name","like","Joh%").Get()
db().Table("users").Where([][]any{
	{"status", "=", 1},
	{"subscribed", "<>", 1},
}).Get()

// Or Where Clauses
db().Table("users").Where("votes",">",100).OrWhere("name", "John").Get()
db().Table("users").Where("votes",">",100).OrWhere(func(query dbgo.SubQuery){
	query.Where("name","Abigail").Where("votes",">",50)
}).Get()
// The example above will produce the following SQL:
// select * from users where votes > 100 or (name = "Abigail" and votes > 50)

// Where Not Clauses
db().Table("products").WhereNot(func(query dbgo.SubQuery){
        query.Where("clearance", true).OrWhere("price", "<", 10)
    }).Get()

// JSON Where Clauses
db().Table("users").Where("preferences->dining->meal", "salad").Get()
db().Table("users").WhereJsonContains("optionsArr->languages", "en").Get()
db().Table("users").WhereJsonContains("optionsArr->languages", []string{"en", "de"}).Get()
db().Table("users").WhereJsonLength("optionsArr->languages", 0).Get()
db().Table("users").WhereJsonLength("optionsArr->languages", ">", 1).Get()

// Additional Where Clauses
db().Table("users").WhereBetween("votes", []int{1, 100}).Get()
db().Table("users").WhereNotBetween("votes", []int{1, 100}).Get()
db().Table("patients").WhereBetweenColumns("weight", []string{"minimum_allowed_weight", "maximum_allowed_weight"}).Get()
db().Table("patients").whereNotBetweenColumns("weight", []string{"minimum_allowed_weight", "maximum_allowed_weight"}).Get()
db().Table("users").WhereIn("id", [1, 2, 3]).Get()
db().Table("users").WhereNotIn("id", [1, 2, 3]).Get()

// sub query
activeUsers := db().Table("users").Select("id").Where("is_active", 1)
db().Table("comments").WhereIn("user_id", activeUsers).Get()
// The example above will produce the following SQL:
// select * from comments where user_id in (
//    select id from users where is_active = 1
//)

// whereNull / whereNotNull / orWhereNull / orWhereNotNull
jc.Ttable("users").WhereNull("updated_at").Get()
jc.Ttable("users").WhereNotNull("updated_at").Get()

// whereDate / whereMonth / whereDay / whereYear / whereTime
jc.Ttable("users").WhereDate("created_at", "2016-12-31").Get()
jc.Ttable("users").WhereMonth("created_at", "12").Get()
jc.Ttable("users").WhereDay("created_at", "31").Get()
jc.Ttable("users").WhereYear("created_at", "2016").Get()
jc.Ttable("users").WhereTime("created_at", "=", "11:20:45").Get()

// whereColumn / orWhereColumn
jc.Ttable("users").WhereColumn("first_name", "last_name").Get()
jc.Ttable("users").WhereColumn("updated_at", ">", "created_at").Get()
db().Table("users").WhereColumn([][]string{
    {"first_name", "=", "last_name"},
    {"updated_at", ">", "created_at"},
}).Get()

// Logical Grouping
db().Table("users").Where("name", "=", "John").Where(func (query dbgo.WhereClause) {
        query.Where("votes", ">", 100).OrWhere("title", "=", "Admin")
    }).Get()
// The example above will produce the following SQL:
// select * from users where name = "John" and (votes > 100 or title = "Admin")
```
### Advanced Where Clauses
```go
// Where Exists Clauses
db().Table("users").WhereExists(func (query dbgo.Database) {
        query.Select(dbgo.Raw(1)).Table("orders").WhereColumn("orders.user_id", "users.id")
    }).Get()

orders := db().Table("orders").Select(dbgo.Raw(1)).WhereColumn("orders.user_id", "users.id")
db().Table("users").WhereExists(orders).Get()
// Both of the examples above will produce the following SQL:
// select * from users where exists (
//    select 1 from orders where orders.user_id = users.id)

// Subquery Where Clauses
db().Table("users").Where(func (query dbgo.SubQuery) {
        query.Select("type").From("membership").WhereColumn("membership.user_id", "users.id").OrderByDesc("membership.start_date").Limit(1)
    }, "=", "Pro").Get()
db().Table("income").Where("amount", "<", func (query dbgo.SubQuery) {
        query.SelectRaw("avg(i.amount)").From("incomes as i")
    }).Get()

// Full Text Where Clauses: match(bio) against("web developer")
db().Table("users").WhereFullText("bio", "web developer").Get()
```
### Ordering, Grouping, Limit and Offset
```go
// Ordering
db().Table("users").OrderBy("name", "desc").Get()
db().Table("users").OrderBy("name", "desc").OrderBy("email", "asc").Get()
db().Table("users").Latest().First()
db().Table("users").InRandomOrder().First()

query := db().Table("users").OrderBy("name")
query.Reorder().Get()

query := db().Table("users").OrderBy("name")
query.Reorder("email", "desc").Get()

// Grouping
db().Table("users").GroupBy("account_id").Having("account_id", ">", 100).Get()
db().Table("orders").SelectRaw("count(id) as number_of_orders, customer_id").GroupBy("customer_id").HavingBetween("number_of_orders", [5, 15]).Get()
db().Table("users").GroupBy("first_name", "status").Having("account_id", ">", 100).Get()

// Limit and Offset and Page
db().Table("users").Skip(10).Take(5).Get()
db().Table("users").Offset(10).Limit(5).Get()
db().Table("users").Page(3).Limit(5).Get()
```
### Conditional Clauses
```go
role := http.Request.Param("role")
db().Table("users").When(role, func (query dbgo.SubQuery, role string) {
        query.Where("role_id", role);
    }).Get()

sortByVotes := http.Request.Param("sort_by_votes").Boolean()
db().Table("users").
    When(sortByVotes, func (query dbgo.SubQuery, sortByVotes bool) {
        query.OrderBy("votes")
    }, func (query dbgo.SubQuery) {
        query.OrderBy("name")
    }).Get()
```
### Insert Statements
```go
db().Table("users").Insert([
    "email" => "kayla@example.com",
    "votes" => 0
]);

db().Table("users").Insert([
    ["email" => "picard@example.com", "votes" => 0],
    ["email" => "janeway@example.com", "votes" => 0],
]);

db().Table("users").InsertOrIgnore([
    ["id" => 1, "email" => "sisko@example.com"],
    ["id" => 2, "email" => "archer@example.com"],
]);

db().Table("pruned_users").InsertUsing([
    "id", "name", "email", "email_verified_at"
], db().Table("users").Select(
    "id", "name", "email", "email_verified_at"
).Where("updated_at", "<=", now().SubMonth()));

// Auto-Incrementing IDs
id := db().Table("users").InsertGetId(
    ["email" => "john@example.com", "votes" => 0]
);

db().Table("flights").Upsert(
    [
        ["departure" => "Oakland", "destination" => "San Diego", "price" => 99],
        ["departure" => "Chicago", "destination" => "New York", "price" => 150]
    ],
    ["departure", "destination"],
    ["price"]
);
```
### Update Statements
```go
affected := db().Table("users").Where("id", 1).Update(["votes" => 1])

db().Table("users").UpdateOrInsert(
    ["email" => "john@example.com", "name" => "John"],
    ["votes" => "2"]
)

affected := db().Table("users").Where("id", 1).Update(["optionsArr->Enabled" => true])

db().Table("users").Increment("votes")

db().Table("users").Increment("votes", 5)

db().Table("users").Decrement("votes")

db().Table("users").Decrement("votes", 5)

db().Table("users").Increment("votes", 1, ["name" => "John"])

db().Table("users").IncrementEach([
    "votes" => 5,
    "balance" => 100,
])
```
### Delete Statements
```go
db().Table('users').Delete()
db().Table('users').Where('votes', '>', 100).Delete()
db().Table('users').Truncate()
```
### Pessimistic Locking
```go
db().Table('users').Where('votes', '>', 100).SharedLock().Get()
db().Table('users').Where('votes', '>', 100).LockForUpdate().Get()
```
### Debugging
```go
db().Table('users').Where('votes', '>', 100).Print()
db().Table('users').Where('votes', '>', 100).PrintRawSql()
```
