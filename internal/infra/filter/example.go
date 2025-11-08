package filter

import "silo/pkg/database/search"

type ExamplePageFilter struct {
	search.Page
	Val  int64  `search:"op:$eq;column:val"`
	Sort string `search:"op:$sort;column:id"`
}
