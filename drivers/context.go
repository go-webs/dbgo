package drivers

type Context struct {
	//handlers []gin.HandlerFunc
	Queries string
	Args    []any
	Err     error
	Prefix  string
}
