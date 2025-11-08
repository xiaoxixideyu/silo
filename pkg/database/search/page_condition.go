package search

// MakePageCondition .
func MakePageCondition(tableName string, pageIndex int, pageSize int, request any) PageCondition {
	condition := &PageConditionImpl{
		pageIndex: pageIndex,
		pageSize:  pageSize,
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

var _ PageCondition = (*PageConditionImpl)(nil)

// PageConditionImpl .
type PageConditionImpl struct {
	ConditionImpl
	pageIndex int
	pageSize  int
}

// GetPageIndex .
func (p *PageConditionImpl) GetPageIndex() int {
	if p.pageIndex <= 0 {
		p.pageIndex = 1
	}
	return p.pageIndex
}

// GetPageSize .
func (p *PageConditionImpl) GetPageSize() int {
	if p.pageSize <= 0 {
		p.pageSize = DefaultMinPageSize
	}
	if p.pageSize > MaxPageSize {
		p.pageSize = MaxPageSize
	}

	return p.pageSize
}

// GetOffset .
func (p *PageConditionImpl) GetOffset() int {
	offset := (p.GetPageIndex() - 1) * p.GetPageSize()
	if offset < 0 {
		offset = 0
	}

	return offset
}
