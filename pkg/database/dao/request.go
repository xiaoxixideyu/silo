package dao

import (
	"silo/pkg/database/interfaces"
	"silo/pkg/database/search"
	"strings"
)

// NewRequest new request
func NewRequest() *Request {
	return &Request{
		filter:  interfaces.D{},
		update:  interfaces.M{},
		project: interfaces.D{},
		sort:    interfaces.D{},
	}
}

// Request condition
type Request struct {
	filter  interfaces.D
	update  interfaces.M
	project interfaces.D
	sort    interfaces.D
	upsert  bool
	limit   int64
	group   []string
	cursor  int
	count   int
}

// GetFilter return filter
func (r *Request) GetFilter() interfaces.D {
	return r.filter
}

// GetLimit return limit
func (r *Request) GetLimit() int64 {
	return r.limit
}

// GetProject return project
func (r *Request) GetProject() interfaces.D {
	return r.project
}

// GetUpdate return update
func (r *Request) GetUpdate() interfaces.M {
	return r.update
}

// GetUpsert return upsert status
func (r *Request) GetUpsert() bool {
	return r.upsert
}

// GetSort sort
func (r *Request) GetSort() interfaces.D {
	return r.sort
}

// GetSort sort
func (r *Request) GetGroup() []string {
	return r.group
}

// FilterD condition by D
func (r *Request) FilterD(f interfaces.D) *Request {
	for _, d := range f {
		r.filter = append(r.filter, d)
	}

	return r
}

// FilterM condition by map
func (r *Request) FilterM(f interfaces.M) *Request {
	for k, v := range f {
		r.filter = append(r.filter, interfaces.E{Key: k, Value: v})
	}

	return r
}

// Filter condition
func (r *Request) Filter(k string, v any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: v})

	return r
}

// FilterIn condition
func (r *Request) FilterIn(k string, v any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPIn: v}})

	return r
}

func (r *Request) FilterNotIn(k string, v any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPNotIn: v}})
	return r
}

// FilterLt condition
func (r *Request) FilterLt(k string, v any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPLT: v}})

	return r
}

// FilterGte .
func (r *Request) FilterGte(k string, v any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPGTE: v}})

	return r
}

// FilterLte .
func (r *Request) FilterLte(k string, v any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPLTE: v}})

	return r
}

// FilterGt condition
func (r *Request) FilterGt(k string, v any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPGT: v}})

	return r
}

// FilterOP condition
func (r *Request) FilterOP(op string, k string, v ...any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{op: v}})

	return r
}

func (r *Request) FilterLike(k string, v any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPLike: v}})

	return r
}

// FilterRange condition
func (r *Request) FilterRange(k string, min, max any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPBetween: []any{min, max}}})

	return r
}

// FilterInM .
func (r *Request) FilterInM(filters []interfaces.M) *Request {
	r.filter = append(r.filter, interfaces.E{Key: search.OPIn, Value: filters})

	return r
}

// FilterIs condition
func (r *Request) FilterIs(k string, v any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPIS: v}})
	return r
}

// FilterIsNotNil condition
func (r *Request) FilterIsNotNil(k string) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPIsNotNil: true}})
	return r
}

// FilterIsNil condition
func (r *Request) FilterIsNil(k string) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPIsNil: true}})
	return r
}

// FilterNe .
func (r *Request) FilterNe(k string, v any) *Request {
	r.filter = append(r.filter, interfaces.E{Key: k, Value: interfaces.M{search.OPNotEqual: v}})
	return r
}

func (r *Request) FilterExpr(expr string) *Request {
	r.filter = append(r.filter, interfaces.E{Key: "", Value: interfaces.M{search.OPExpr: expr}})
	return r
}

// SetLimit limit
func (r *Request) SetLimit(limit int64) *Request {
	r.limit = limit
	return r
}

// Inc a filed
func (r *Request) Inc(k string, v any) *Request {
	_, ok := r.update[search.OPInc]
	if !ok {
		r.update[search.OPInc] = interfaces.M{}
	}
	m := r.update[search.OPInc].(interfaces.M)
	m[k] = v

	return r
}

// Set a filed
func (r *Request) Set(k string, v any) *Request {
	_, ok := r.update[search.OPSet]
	if !ok {
		r.update[search.OPSet] = interfaces.M{}
	}
	m := r.update[search.OPSet].(interfaces.M)
	m[k] = v

	return r
}

// SetOnInsert a filed
func (r *Request) SetOnInsert(k string, v any) *Request {
	_, ok := r.update[search.OPSetOnInsert]
	if !ok {
		r.upsert = true
		r.update[search.OPSetOnInsert] = interfaces.M{}
	}
	m := r.update[search.OPSetOnInsert].(interfaces.M)
	m[k] = v
	return r
}

// SetOnInsertM a filed
func (r *Request) SetOnInsertM(insert interfaces.M) *Request {
	_, ok := r.update[search.OPSetOnInsert]
	if !ok {
		r.upsert = true
		r.update[search.OPSetOnInsert] = interfaces.M{}
	}
	m := r.update[search.OPSetOnInsert].(interfaces.M)
	for k, v := range insert {
		m[k] = v
	}

	return r
}

// Upsert mark upsert flag
func (r *Request) Upsert() *Request {
	_, ok := r.update[search.OPSetOnInsert]
	if !ok {
		r.upsert = true
		r.update[search.OPSetOnInsert] = interfaces.M{}
	}

	return r
}

// Project include a field
func (r *Request) Project(keys ...string) *Request {
	for _, k := range keys {
		r.project = append(r.project, interfaces.E{Key: k, Value: 1})
	}

	return r
}

// UnProject Exclude a field
func (r *Request) UnProject(keys ...string) *Request {
	for _, k := range keys {
		r.project = append(r.project, interfaces.E{Key: k, Value: 0})
	}

	return r
}

// Sort v should be asc or desc
func (r *Request) Sort(k string, v string) *Request {
	switch strings.ToUpper(v) {
	case search.SortASC:
		r.SortASC(k)
	case search.SortDESC:
		r.SortDESC(k)
	}

	return r
}

// SortASC sort by asc
func (r *Request) SortASC(k string) *Request {
	if r.sort == nil {
		r.sort = interfaces.D{}
	}

	r.sort = append(r.sort, interfaces.E{Key: k, Value: interfaces.ASC})

	return r
}

// SortDESC sort by desc
func (r *Request) SortDESC(k string) *Request {
	if r.sort == nil {
		r.sort = interfaces.D{}
	}

	r.sort = append(r.sort, interfaces.E{Key: k, Value: interfaces.DESC})

	return r
}

func (r *Request) Group(k string) *Request {
	if r.group == nil {
		r.group = []string{}
	}
	r.group = append(r.group, k)
	return r
}

func (r *Request) Select(keys ...string) *Request {
	for _, k := range keys {
		r.project = append(r.project, interfaces.E{Key: k, Value: 1})
	}
	return r
}

func (r *Request) SetCursor(cursor int) *Request {
	r.cursor = cursor
	return r
}

func (r *Request) SetCount(count int) *Request {
	r.count = count
	return r
}

func (r *Request) GetCursor() int {
	return r.cursor
}

func (r *Request) GetCount() int {
	if r.count <= 0 {
		return 10
	}
	return r.count
}
