package search

// MakeCursorCondition 创建游标条件
func MakeCursorCondition(tableName string, cursor int, count int, request any) CursorCondition {
	condition := &CursorConditionImpl{
		cursor: cursor,
		count:  count,
	}

	ParseSearchTag(tableName, request, condition)

	// Projection
	if p, ok := request.(Projection); ok {
		condition.projection = p.GetProjection()
	}

	// Distinct
	if p, ok := request.(Distinct); ok {
		condition.distinct = p.GetDistinct()
	}

	return condition
}

var _ CursorCondition = (*CursorConditionImpl)(nil)

// CursorConditionImpl 游标条件实现
type CursorConditionImpl struct {
	ConditionImpl
	cursor int
	count  int
}

// GetCursor 获取游标位置
func (c *CursorConditionImpl) GetCursor() int {
	if c.cursor < 0 {
		c.cursor = 0
	}
	return c.cursor
}

// GetCount 获取查询数量
func (c *CursorConditionImpl) GetCount() int {
	if c.count <= 0 {
		c.count = DefaultMinPageSize
	}
	if c.count > MaxPageSize {
		c.count = MaxPageSize
	}
	return c.count
}
