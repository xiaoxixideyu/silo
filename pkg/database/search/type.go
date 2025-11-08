package search

// QueryTag .
const QueryTag = "search"

// Projection .
type Projection interface {
	GetProjection() []string
}

// Distinct .
type Distinct interface {
	GetDistinct() []string
}

// Condition .
type Condition interface {
	Projection
	Distinct
	AddSearchTag(tag *Tag)
	GetSearchTags() []*Tag
}

// Pagination .
type Pagination interface {
	GetPageIndex() int
	GetPageSize() int
	GetOffset() int
}

// PageCondition .
type PageCondition interface {
	Condition
	Pagination
}

// Cursor 游标查询接口
type Cursor interface {
	GetCursor() int
	GetCount() int
}

// CursorCondition 游标条件接口
type CursorCondition interface {
	Condition
	Cursor
}
