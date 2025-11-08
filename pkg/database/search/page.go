package search

var _ Pagination = (*Page)(nil)

// Page .
type Page struct {
	PageIndex int `json:"page_index" query:"pageIndex"`
	PageSize  int `json:"page_size" query:"pageSize"`
}

// GetPageIndex .
func (p *Page) GetPageIndex() int {
	if p.PageIndex <= 0 {
		p.PageIndex = 1
	}
	return p.PageIndex
}

// GetPageSize .
func (p *Page) GetPageSize() int {
	if p.PageSize <= 0 {
		p.PageSize = DefaultMinPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
	return p.PageSize
}

// GetOffset .
func (p *Page) GetOffset() int {
	offset := (p.GetPageIndex() - 1) * p.GetPageSize()
	if offset < 0 {
		offset = 0
	}

	return offset
}

// NewPageResult .
func NewPageResult(pageIndex int, pageSize int, totalCount int64) *PageResult {
	return &PageResult{
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: totalCount,
	}
}

// PageResult .
type PageResult struct {
	PageIndex  int   `json:"page_index"`
	PageSize   int   `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	// Results    []any `json:"results"`
}

// SetPageIndex .
func (p *PageResult) SetPageIndex(index int) *PageResult {
	p.PageIndex = index
	return p
}

// SetPageSize .
func (p *PageResult) SetPageSize(size int) *PageResult {
	p.PageSize = size
	return p
}

// SetTotalCount .
func (p *PageResult) SetTotalCount(count int64) *PageResult {
	p.TotalCount = count
	return p
}

// // SetResults .
// func (p *PageResult) SetResults(results []any) *PageResult {
// 	p.Results = results
// 	return p
// }
