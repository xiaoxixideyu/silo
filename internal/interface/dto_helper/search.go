package dto_helper

import "silo/pkg/database/search"

// Page .
func Page(in search.Page) search.Page {
	return in
}

// PagePoint .
func PagePoint(in *search.Page) *search.Page {
	return in
}

// PageResult .
func PageResult(in search.PageResult) search.PageResult {
	return in
}

// PageResultPoint .
func PageResultPoint(in *search.PageResult) *search.PageResult {
	return in
}
