package iface

type Context struct {
	//handlers []gin.HandlerFunc
	Query string
	Args  []any
	Err   error
}
