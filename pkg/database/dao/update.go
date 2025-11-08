package dao

import (
	"context"
	"fmt"
	"silo/pkg/database/interfaces"
	"silo/pkg/database/seedwork"
	"silo/pkg/errs"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpdateMany .
func (r *DAOImpl) UpdateMany(ctx context.Context, model interfaces.Model, req interfaces.Update) error {
	db := r.DB(ctx).Table(model.TableName()).Model(model)
	db = parseFilter(db, req)
	db = parseUpdate(db, req.GetUpdate())

	return ErrorMap(db.Error, model.TableName())
}

// FindOneAndUpdate .
func (r *DAOImpl) FindOneAndUpdate(ctx context.Context, tableName string, req interfaces.Update, res any) error {
	db := r.DB(ctx).Table(tableName).Model(res).Clauses(clause.Returning{})
	db = parseFilter(db, req)
	db = parseUpdate(db, req.GetUpdate())

	return ErrorMap(db.Error, tableName)
}

// Save .
func (r *DAOImpl) Save(ctx context.Context, obj interfaces.Model) error {
	changes := obj.GetChanges()
	if len(changes) == 0 {
		return nil
	}

	if w, ok := obj.(interfaces.WithTimestamp); ok {
		w.WithTimestampSetUpdatedAt()
		changes[seedwork.TableColumnUpdatedAt] = w.GetUpdatedAt()
	}

	var optimistic = false
	var where map[string]any
	if w, ok := obj.(interfaces.WithVersion); ok {
		currentVersion := w.GetVersion()
		w.IncVersion()
		changes[seedwork.TableColumnVersion] = w.GetVersion()

		optimistic = true
		where = map[string]any{
			seedwork.TableColumnVersion: currentVersion,
		}
	}

	// Model(obj) include primary key
	db := r.DB(ctx).Model(obj)
	if optimistic {
		db.Where(where)
	}

	res := db.Updates(changes)
	if res.Error != nil {
		return ErrorMap(res.Error, obj.TableName())
	}

	if optimistic && res.RowsAffected == 0 {
		return errs.NewCustomError(errs.StatusOptimisticLockError, "Optimistic lock failed").
			WithError(fmt.Errorf("tableName: %s, idMap: %v, where: %v", obj.TableName(), obj.IDMap(), where))
	}
	return nil
}

const defaultBatchSize = 100

// SaveInBatches 批量保存记录，支持乐观锁和非乐观锁场景
func (r *DAOImpl) SaveInBatches(ctx context.Context, models []interfaces.Model) error {
	if len(models) == 0 {
		return nil
	}

	// 将记录分为乐观锁和非乐观锁两组
	var optimisticModels, normalModels []interfaces.Model
	for _, model := range models {
		if _, ok := model.(interfaces.WithVersion); ok {
			optimisticModels = append(optimisticModels, model)
		} else {
			normalModels = append(normalModels, model)
		}
	}

	// 处理非乐观锁记录
	for i := 0; i < len(normalModels); i += defaultBatchSize {
		end := i + defaultBatchSize
		if end > len(normalModels) {
			end = len(normalModels)
		}
		batch := normalModels[i:end]

		if err := r.DB(ctx).Transaction(func(tx *gorm.DB) error {
			for _, model := range batch {
				changes := model.GetChanges()
				if len(changes) == 0 {
					continue
				}

				if w, ok := model.(interfaces.WithTimestamp); ok {
					w.WithTimestampSetUpdatedAt()
					changes[seedwork.TableColumnUpdatedAt] = w.GetUpdatedAt()
				}

				if err := tx.Model(model).Where(model.IDMap()).Updates(changes).Error; err != nil {
					return ErrorMap(err, model.TableName())
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	// 处理乐观锁记录
	for i := 0; i < len(optimisticModels); i += defaultBatchSize {
		end := i + defaultBatchSize
		if end > len(optimisticModels) {
			end = len(optimisticModels)
		}
		batch := optimisticModels[i:end]

		if err := r.DB(ctx).Transaction(func(tx *gorm.DB) error {
			for _, model := range batch {
				changes := model.GetChanges()
				if len(changes) == 0 {
					continue
				}

				if w, ok := model.(interfaces.WithTimestamp); ok {
					w.WithTimestampSetUpdatedAt()
					changes[seedwork.TableColumnUpdatedAt] = w.GetUpdatedAt()
				}

				w := model.(interfaces.WithVersion)
				currentVersion := w.GetVersion()
				w.IncVersion()
				changes[seedwork.TableColumnVersion] = w.GetVersion()

				result := tx.Model(model).
					Where(model.IDMap()).
					Where(map[string]any{seedwork.TableColumnVersion: currentVersion}).
					Updates(changes)

				if result.Error != nil {
					return ErrorMap(result.Error, model.TableName())
				}

				if result.RowsAffected == 0 {
					return errs.NewCustomError(errs.StatusOptimisticLockError, "Optimistic lock failed").
						WithError(fmt.Errorf("tableName: %s, idMap: %v", model.TableName(), model.IDMap()))
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
