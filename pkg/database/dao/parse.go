package dao

import (
	"fmt"
	"silo/pkg/database/interfaces"
	"silo/pkg/database/search"
	"silo/pkg/utils"

	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

func parseFilter(db *gorm.DB, req interfaces.Filter) *gorm.DB {
	if project := req.GetProject(); len(project) > 0 {
		cols := []string{}
		for _, v := range project {
			if v.Value == 1 {
				cols = append(cols, v.Key)
			}
		}
		if len(cols) > 0 {
			db = db.Select(cols)
		}
	}

	filter := req.GetFilter()
	for _, v := range filter {
		if m, ok := v.Value.(interfaces.M); ok {
			db = parseFilterOP(db, v.Key, m)
		} else {
			db = db.Where(fmt.Sprintf("%s = ?", v.Key), v.Value)
		}
	}

	if limit := req.GetLimit(); limit > 0 {
		db = db.Limit(int(limit))
	}
	sort := req.GetSort()
	for _, v := range sort {
		if v.Value == interfaces.ASC {
			db = db.Order(fmt.Sprintf("%s ASC", v.Key))
		} else if v.Value == interfaces.DESC {
			db = db.Order(fmt.Sprintf("%s DESC", v.Key))
		}
	}

	for _, v := range req.GetGroup() {
		db = db.Group(v)
	}

	return db
}

func parseFilterOP(db *gorm.DB, column string, m interfaces.M) *gorm.DB {
	for k, v := range m {
		switch k {
		case search.OPEqual:
			db = db.Where(fmt.Sprintf("%s = ?", column), v)
		case search.OPEqualIgnoreCase:
			db = db.Where(fmt.Sprintf("UPPER(%s) = UPPER(?)", column), v)
		case search.OPNotEqual:
			db = db.Where(fmt.Sprintf("%s != ?", column), v)
		case search.OPGT:
			db = db.Where(fmt.Sprintf("%s > ?", column), v)
		case search.OPGTE:
			db = db.Where(fmt.Sprintf("%s >= ?", column), v)
		case search.OPLT:
			db = db.Where(fmt.Sprintf("%s < ?", column), v)
		case search.OPLTE:
			db = db.Where(fmt.Sprintf("%s <= ?", column), v)
		case search.OPIS:
			db = db.Where(column, v)
		case search.OPIn:
			db = db.Where(fmt.Sprintf("%s IN (?)", column), v)
		case search.OPNotIn:
			db = db.Where(fmt.Sprintf("%s NOT IN (?)", column), v)
		case search.OPBetween:
			db = db.Where(fmt.Sprintf("%s BETWEEN (?)", column), v)
		case search.OPLike:
			db = db.Where(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%v%%", v))
		case search.OPContains:
			str, _ := utils.ToJSON(v)
			db = db.Where(fmt.Sprintf("JSON_CONTAINS(%s,?)", column), str)
		case search.OPArrayContains:
			db = db.Where(fmt.Sprintf("%s @> (?)", column), v)
		case search.OPExpr:
			db = db.Where(fmt.Sprintf("(%v)", v))
		case search.OPIsNotNil:
			db = db.Where(fmt.Sprintf("%s IS NOT NULL", column))
		case search.OPIsNil:
			db = db.Where(fmt.Sprintf("%s IS NULL", column))
		}
	}

	return db
}

func parseUpdate(db *gorm.DB, m interfaces.M) *gorm.DB {
	updates := interfaces.M{}
	for k, v := range m {
		switch k {
		case search.OPSet:
			if set, ok := v.(interfaces.M); ok {
				maps.Copy(updates, set)
			}
		case search.OPInc:
			if set, ok := v.(interfaces.M); ok {
				for col, inc := range set {
					updates[col] = gorm.Expr(fmt.Sprintf("%s + ?", col), inc)
				}
			}
		}
	}

	if db.Statement.Schema == nil {
		db = ensureSchema(db)
	}

	return db.Updates(updates)
}

func ensureSchema(db *gorm.DB) *gorm.DB {
	if db.Statement == nil {
		return db
	}
	if db.Statement.Model == nil && db.Statement.Dest != nil {
		db = db.Model(db.Statement.Dest)
	}
	if db.Statement.Schema == nil {
		_ = db.Statement.Parse(db.Statement.Model)
	}
	return db
}
