package dao

import (
	"context"
	"silo/pkg/database/interfaces"
)

// Create record
func (r *DAOImpl) Create(ctx context.Context, tableName string, model interfaces.Model) error {
	res := r.DB(ctx).Table(tableName).Create(model)
	if res.Error != nil {
		return ErrorMap(res.Error, tableName)
	}

	return nil
}

// CreateMany records
func (r *DAOImpl) CreateMany(ctx context.Context, tableName string, models any) error {
	res := r.DB(ctx).Table(tableName).Create(models)
	if res.Error != nil {
		return ErrorMap(res.Error, tableName)
	}

	return nil
}

// FindOrCreate returns the record if it exists, otherwise creates it
func (r *DAOImpl) FindOrCreate(ctx context.Context, tableName string, model interfaces.Model, res any) (create bool, err error) {
	result := r.DB(ctx).Where(model.IDMap()).Attrs(model).FirstOrCreate(res)
	if result.Error != nil {
		return false, ErrorMap(result.Error, tableName)
	}

	return result.RowsAffected > 0, nil
}

// CreateInBatches .
func (r *DAOImpl) CreateInBatches(ctx context.Context, tableName string, models any, batchSize int) error {
	res := r.DB(ctx).Table(tableName).CreateInBatches(models, batchSize)
	if res.Error != nil {
		return ErrorMap(res.Error, tableName)
	}

	return nil
}
