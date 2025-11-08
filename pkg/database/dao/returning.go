package dao

import (
	"context"
	"silo/pkg/database/interfaces"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IncreaseResult 包含原始值和更新后的值
type IncreaseResult struct {
	OriginalValue int64 `json:"original_value"`
	NewValue      int64 `json:"new_value"`
}

// IncreaseAndReturn .
func (r *DAOImpl) IncreaseAndReturn(ctx context.Context, tableName string, req interfaces.Filter, column string, step int64, res any) error {
	db := r.DB(ctx).Clauses(clause.Returning{}).Table(tableName).Model(res)
	db = parseFilter(db, req)
	db = db.Update(column, gorm.Expr(column+" + ?", step))

	return db.Error
}

// IncreaseAndReturnValues 返回原始值和更新后的自增值
func (r *DAOImpl) IncreaseAndReturnValues(ctx context.Context, tableName string, req interfaces.Filter, column string, step int64, res any) (*IncreaseResult, error) {
	// 使用事务确保数据一致性
	return r.performIncreaseWithValues(ctx, tableName, req, column, step, res)
}

// performIncreaseWithValues 在事务中执行查询和更新操作
func (r *DAOImpl) performIncreaseWithValues(ctx context.Context, tableName string, req interfaces.Filter, column string, step int64, res any) (*IncreaseResult, error) {
	var result IncreaseResult

	err := r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 先查询原始值
		var originalRecord map[string]any
		queryDB := tx.Table(tableName)
		queryDB = parseFilter(queryDB, req)

		if err := queryDB.First(&originalRecord).Error; err != nil {
			return ErrorMap(err, tableName)
		}

		// 获取原始值
		if val, ok := originalRecord[column]; ok {
			switch v := val.(type) {
			case int64:
				result.OriginalValue = v
			case int32:
				result.OriginalValue = int64(v)
			case int:
				result.OriginalValue = int64(v)
			case float64:
				result.OriginalValue = int64(v)
			default:
				result.OriginalValue = 0
			}
		}

		// 计算新值
		result.NewValue = result.OriginalValue + step

		// 2. 执行更新操作
		updateDB := tx.Clauses(clause.Returning{}).Table(tableName).Model(res)
		updateDB = parseFilter(updateDB, req)
		updateDB = updateDB.Update(column, gorm.Expr(column+" + ?", step))

		return updateDB.Error
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}
