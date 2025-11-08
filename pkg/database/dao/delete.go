package dao

import (
	"context"
	"silo/pkg/database/interfaces"
)

// Delete .
func (r *DAOImpl) Delete(ctx context.Context, tableName string, model interfaces.Model) error {
	if err := r.DB(ctx).Table(tableName).Where(model.IDMap()).Delete(model).Error; err != nil {
		return ErrorMap(err, tableName)
	}
	return nil
}

// DeleteOne .
func (r *DAOImpl) DeleteOne(ctx context.Context, tableName string, req interfaces.Filter) error {
	db := r.DB(ctx).Table(tableName)
	db = parseFilter(db, req)
	result := db.Delete(map[string]any{})
	return ErrorMap(result.Error, tableName)
}

// DeleteMany .
func (r *DAOImpl) DeleteMany(ctx context.Context, tableName string, req interfaces.Filter) error {
	db := r.DB(ctx).Table(tableName)
	db = parseFilter(db, req)
	result := db.Delete(map[string]any{})
	return ErrorMap(result.Error, tableName)
}
