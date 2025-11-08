package dao

import (
	"fmt"
	"silo/pkg/database/search"
	"silo/pkg/utils"
	"strings"

	"github.com/samber/lo"
)

// DBCondition .
type DBCondition interface {
	SetWhere(k string, v ...any)
	SetOr(k string, v ...any)
	SetOrder(k string)
	SetJoin(on string) DBCondition
	SetGroup(k string)
	SetGroupOr() DBCondition
}

// NewDBCondition .
func NewDBCondition(q search.Condition) *GormCondition {
	tags := q.GetSearchTags()
	cond := &GormCondition{}

	resolveDBCondition(tags, cond, false)
	return cond
}

// GormCondition .
type GormCondition struct {
	GormBase
	GroupOr *GormCondition
	Join    []*GormJoin
}

// SetJoin .
func (c *GormCondition) SetJoin(on string) DBCondition {
	if c.Join == nil {
		c.Join = make([]*GormJoin, 0)
	}
	join := &GormJoin{
		On:       on,
		GormBase: GormBase{},
	}
	c.Join = append(c.Join, join)

	return join
}

// SetGroupOr .
func (c *GormCondition) SetGroupOr() DBCondition {
	if c.GroupOr == nil {
		c.GroupOr = &GormCondition{}
	}

	return c.GroupOr
}

// GormJoin .
type GormJoin struct {
	GormBase
	On string
}

// SetJoin .
func (e *GormJoin) SetJoin(on string) DBCondition {
	return nil
}

// GormBase .
type GormBase struct {
	Where   map[string][]any
	Order   []string
	Group   []string
	Or      map[string][]any
	GroupOr *GormBase
}

// SetWhere condition
func (c *GormBase) SetWhere(k string, v ...any) {
	if c.Where == nil {
		c.Where = make(map[string][]any)
	}

	c.Where[k] = v
}

// SetOr .
func (c *GormBase) SetOr(k string, v ...any) {
	if c.Or == nil {
		c.Or = make(map[string][]any)
	}

	c.Or[k] = v
}

// SetOrder .
func (c *GormBase) SetOrder(k string) {
	if c.Order == nil {
		c.Order = make([]string, 0)
	}

	c.Order = append(c.Order, k)
}

// SetGroup .
func (c *GormBase) SetGroup(k string) {
	if c.Group == nil {
		c.Group = make([]string, 0)
	}

	c.Group = append(c.Group, k)
}

// SetGroupOr .
func (c *GormBase) SetGroupOr() DBCondition {
	return nil
}

func resolveDBCondition(tags []*search.Tag, cond DBCondition, or bool) {
	f := cond.SetWhere
	if or {
		f = cond.SetOr
	}
	for _, t := range tags {
		switch t.Operator {
		case search.OPLeftJoin:
			joinStr := parseJoin("LEFT JOIN", t)
			join := cond.SetJoin(joinStr)
			if len(t.Children) > 0 {
				resolveDBCondition(t.Children, join, false)
			}
		case search.OPJoin:
			joinStr := parseJoin("INNER JOIN", t)
			join := cond.SetJoin(joinStr)
			if len(t.Children) > 0 {
				resolveDBCondition(t.Children, join, false)
			}
		case search.OPGroupOr:
			groupOr := cond.SetGroupOr()
			if len(t.Children) > 0 {
				resolveDBCondition(t.Children, groupOr, true)
			}
		case search.OPEqual:
			f(fmt.Sprintf("%s = ?", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPEqualIgnoreCase:
			f(fmt.Sprintf("UPPER(%s) = UPPER(?)", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPNotEqual:
			f(fmt.Sprintf("%s != ?", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPGT:
			f(fmt.Sprintf("%s > ?", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPGTE:
			f(fmt.Sprintf("%s >= ?", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPLT:
			f(fmt.Sprintf("%s < ?", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPLTE:
			f(fmt.Sprintf("%s <= ?", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPIS:
			f(EncloseTableColumn(t.Table, t.Column), t.Value)
		case search.OPIsNotNil:
			f(fmt.Sprintf("%s IS NOT NULL", EncloseTableColumn(t.Table, t.Column)))
		case search.OPIsNil:
			f(fmt.Sprintf("%s IS NULL", EncloseTableColumn(t.Table, t.Column)))
		case search.OPIn:
			f(fmt.Sprintf("%s IN (?)", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPNotIn:
			f(fmt.Sprintf("%s NOT IN (?)", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPBetween:
			f(fmt.Sprintf("%s BETWEEN (?)", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPLike:
			f(fmt.Sprintf("%s LIKE ?", EncloseTableColumn(t.Table, t.Column)), fmt.Sprintf("%%%v%%", t.Value))
		case search.OPContains:
			str, _ := utils.ToJSON(t.Value)
			f(fmt.Sprintf("JSON_CONTAINS(%s,?)", EncloseTableColumn(t.Table, t.Column)), str)
		case search.OPPattern:
			f(t.OnRaw)
		case search.OPArrayContains:
			f(fmt.Sprintf("ARRAY_INCLUDES(%s,?)", EncloseTableColumn(t.Table, t.Column)), t.Value)
		case search.OPJsonToStringLike:
			f(fmt.Sprintf("TO_JSON_STRING(%s) LIKE ?", EncloseTableColumn(t.Table, t.Column)), fmt.Sprintf("%%%v%%", t.Value))
		case search.OPJson:
			f(fmt.Sprintf(`JSON_VALUE(%s,"$.%s") = ?`, EncloseTableJson(t.Table, t.Column), t.JsonValue), t.Value)
		case search.OPJsonLike:
			f(fmt.Sprintf(`JSON_VALUE(%s,"$.%s") LIKE ?`, EncloseTableJson(t.Table, t.Column), t.JsonValue), fmt.Sprintf("%%%v%%", t.Value))
		case search.OPSort:
			sortType, ok := t.Value.(string)
			if ok {
				switch strings.ToUpper(sortType) {
				case search.SortASC, search.SortDESC, search.SortDESCNullsLast:
					value := parseSort(t.Table, t.Column, sortType)
					cond.SetOrder(value)
				}
			}
		case search.OPGroup:
			if t.Column != "" {
				cond.SetGroup(t.Column)
			}
		case search.OPCustomSort:
			if customSort, ok := t.Value.(string); ok {
				cond.SetOrder(customSort)
			}
		}
	}
}

func parseJoin(joinType string, t *search.Tag) string {
	as := lo.Ternary(t.As == "", t.Join, t.As)
	on := t.OnRaw
	if on == "" {
		conds := []string{}
		for _, v := range t.On {
			conds = append(conds, fmt.Sprintf(
				"%s %s %s",
				EncloseTableColumn(t.Table, v[0]),
				v[1],
				EncloseTableColumn(as, v[2]),
			))
		}

		on = strings.Join(conds, " AND ")
	}

	return fmt.Sprintf("%s %s AS %s ON %s", joinType, Enclose(t.Join), as, on)
}

func parseSort(tableName string, column string, sort string) string {
	cols := strings.Split(column, ",")

	var elems []string
	for _, col := range cols {
		if tableName == "" {
			elems = append(elems, fmt.Sprintf("%s %s", col, sort))
		} else {
			elems = append(elems, fmt.Sprintf("%s %s", EncloseTableColumn(tableName, col), sort))
		}
	}
	return strings.Join(elems, ", ")
}
