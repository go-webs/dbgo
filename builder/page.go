package builder

type PageBuilder struct {
	limit, offset, page int
}

func NewPageBuilder() *PageBuilder {
	return &PageBuilder{}
}

func (b *PageBuilder) Limit(arg int) *PageBuilder {
	b.limit = arg
	return b
}
func (b *PageBuilder) Page(arg int) *PageBuilder {
	b.page = arg
	return b
}
func (b *PageBuilder) Offset(arg int) *PageBuilder {
	b.offset = arg
	return b
}

func (b *PageBuilder) BuildLimit() (sqlSegment string, binds []any) {
	return
}

func (b *PageBuilder) BuildPage() (sqlSegment string, binds []any) {
	var offset int
	if b.offset > 0 {
		offset = b.offset
	} else if b.page > 0 {
		offset = b.limit * (b.page - 1)
	}
	if b.limit > 0 {
		if offset > 0 {
			sqlSegment = "LIMIT ? OFFSET ?"
			binds = append(binds, b.limit, offset)
		} else {
			sqlSegment = "LIMIT ?"
			binds = append(binds, b.limit)
		}
	}
	return
}
func (b *PageBuilder) GetPagination() (limit, offset, page int) {
	return b.limit, b.offset, b.page
}
