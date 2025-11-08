package search

// MakeCondition .
func MakeCondition(tableName string, request any) Condition {
	condition := &ConditionImpl{}

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

var _ Condition = (*ConditionImpl)(nil)

// ConditionImpl .
type ConditionImpl struct {
	tags       []*Tag
	projection []string
	distinct   []string
}

// AddSearchTag .
func (p *ConditionImpl) AddSearchTag(tag *Tag) {
	p.tags = append(p.tags, tag)
}

// GetSearchTags .
func (p *ConditionImpl) GetSearchTags() []*Tag {
	return p.tags
}

// GetProjection .
func (p *ConditionImpl) GetProjection() []string {
	return p.projection
}

// GetDistinct .
func (p *ConditionImpl) GetDistinct() []string {
	return p.distinct
}
