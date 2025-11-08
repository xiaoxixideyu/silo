package search

var _ Cursor = (*CursorPage)(nil)

// CursorPage 游标分页
type CursorPage struct {
	Cursor int `json:"cursor" query:"cursor"`
	Count  int `json:"count" query:"count"`
}

// GetCursor 获取游标位置
func (c *CursorPage) GetCursor() int {
	if c.Cursor < 0 {
		c.Cursor = 0
	}
	return c.Cursor
}

// GetCount 获取查询数量
func (c *CursorPage) GetCount() int {
	if c.Count <= 0 {
		c.Count = DefaultMinPageSize
	}
	if c.Count > MaxPageSize {
		c.Count = MaxPageSize
	}
	return c.Count
}

func NewCursor(cursor, count int) *CursorPage {
	return &CursorPage{
		Cursor: cursor,
		Count:  count,
	}
}

type CursorResult struct {
	NextCursor int  `json:"next_cursor"`
	HasMore    bool `json:"has_more"`
}

func NewCursorResult(nextCursor int, hasMore bool) *CursorResult {
	return &CursorResult{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}
}
