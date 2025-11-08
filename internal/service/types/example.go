package types

import (
	"silo/internal/infra/entity"
	"silo/pkg/database/search"
)

type ExamplePage struct {
	search.Page
	Val int64
}

type ExamplePageResult struct {
	*search.PageResult
	Infos []*entity.Example
}
