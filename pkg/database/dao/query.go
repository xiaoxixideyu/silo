package dao

import (
	"context"
	"fmt"
	"reflect"
	"silo/pkg/database/interfaces"
	"silo/pkg/database/search"
	"silo/pkg/errs"
	"silo/pkg/tracer"

	"gorm.io/gorm"
)

var _ interfaces.Dao = (*DAOImpl)(nil)
var _ interfaces.TXBeginner = (*DAOImpl)(nil)

// DB DB
func (r *DAOImpl) DB(ctx context.Context) *gorm.DB {
	db := CtxToDB(ctx)
	if db != nil {
		return db
	}

	return r.db.WithContext(ctx)
}

// Transaction .
func (r *DAOImpl) Transaction(txCtx interfaces.TXContext, fn interfaces.TXFunc) error {
	db := r.DB(txCtx)

	tx := tracer.FromContext(txCtx)
	if tx != nil {
		span := tx.StartSpan("tx", "dao", nil)

		defer span.End()
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		c := DBToCtx(txCtx, tx)
		err := fn(c, txCtx)
		return err
	})
	return ErrorMap(err, "tx")
}

// FindOne .
func (r *DAOImpl) FindOne(ctx context.Context, tableName string, req interfaces.Filter, res any) error {
	valueOf := reflect.ValueOf(res)
	if valueOf.Kind() != reflect.Pointer {
		return errs.NewInternalServerError("find res param is not pointer")
	}

	db := r.DB(ctx).Table(tableName)
	db = parseFilter(db, req)

	err := db.First(res).Error

	return ErrorMap(err, tableName)
}

// Find .
func (r *DAOImpl) Find(ctx context.Context, tableName string, req interfaces.Filter, res any) error {
	valueOf := reflect.ValueOf(res)
	if valueOf.Kind() != reflect.Pointer {
		return errs.NewInternalServerError("find res param is not pointer")
	}
	db := r.DB(ctx).Table(tableName)
	db = parseFilter(db, req)

	result := db.Find(res)
	err := result.Error

	return ErrorMap(err, tableName)
}

// FindOneByCondition .
func (r *DAOImpl) FindOneByCondition(ctx context.Context, tableName string, req any, res any) error {
	condition := search.MakeCondition(tableName, req)

	db := r.DB(ctx)
	db = r.applyCondition(db, condition, false)
	err := db.
		Table(tableName).
		First(res).Error

	return ErrorMap(err, tableName)
}

// FindByCondition .
func (r *DAOImpl) FindByCondition(ctx context.Context, tableName string, req any, res any) error {
	condition := search.MakeCondition(tableName, req)

	db := r.DB(ctx)
	db = r.applyCondition(db, condition, false)
	err := db.
		Table(tableName).
		Find(res).Error
	return ErrorMap(err, tableName)
}

// Page .
func (r *DAOImpl) Page(ctx context.Context, tableName string, in search.Pagination, out any) (*search.PageResult, error) {
	pageIndex := in.GetPageIndex()
	pageSize := in.GetPageSize()
	condition := search.MakePageCondition(tableName, pageIndex, pageSize, in)

	// ress
	db := r.DB(ctx)
	db = r.applyCondition(db, condition, false)
	err := db.
		Table(tableName).
		Offset(condition.GetOffset()).
		Limit(condition.GetPageSize()).
		Find(out).
		Error

	if err != nil {
		return nil, ErrorMap(err, tableName)
	}

	// totalCount
	var totalCount int64
	db = r.DB(ctx)
	db = r.applyCondition(db, condition, true)
	err = db.
		Table(tableName).
		Count(&totalCount).
		Error
	if err != nil {
		return nil, ErrorMap(err, tableName)
	}

	pageResult := search.NewPageResult(condition.GetPageIndex(), condition.GetPageSize(), totalCount)

	return pageResult, ErrorMap(err, tableName)
}

// FindInBatches .
func (r *DAOImpl) FindInBatches(ctx context.Context, tableName string, req interfaces.Filter, res any, batchSize int, fn func(batch int) error) error {
	db := r.DB(ctx).Table(tableName)
	db = parseFilter(db, req)

	result := db.FindInBatches(res, batchSize, func(tx *gorm.DB, batch int) error {
		return fn(batch)
	})
	err := result.Error

	return ErrorMap(err, tableName)
}

// FindInBatchesByCondition .
func (r *DAOImpl) FindInBatchesByCondition(ctx context.Context, tableName string, req any, res any, batchSize int, fn func(batch int) error) error {
	condition := search.MakeCondition(tableName, req)

	db := r.DB(ctx).Table(tableName)
	db = r.applyCondition(db, condition, false)

	result := db.FindInBatches(res, batchSize, func(tx *gorm.DB, batch int) error {
		return fn(batch)
	})
	err := result.Error

	return ErrorMap(err, tableName)
}

// Count .
func (r *DAOImpl) Count(ctx context.Context, tableName string, req interfaces.Filter) (int64, error) {
	db := r.DB(ctx).Table(tableName)
	db = parseFilter(db, req)

	var count int64
	err := db.Count(&count).Error

	return count, ErrorMap(err, tableName)
}

func (r *DAOImpl) CountByCondition(ctx context.Context, tableName string, req any) (int64, error) {
	condition := search.MakeCondition(tableName, req)

	// count
	var count int64
	db := r.DB(ctx)
	db = r.applyCondition(db, condition, true)
	err := db.
		Table(tableName).
		Count(&count).
		Error

	return count, ErrorMap(err, tableName)
}

func (r *DAOImpl) Cursor(ctx context.Context, tableName string, req interfaces.Filter, res any) (*search.CursorResult, error) {
	valueOf := reflect.ValueOf(res)
	if valueOf.Kind() != reflect.Pointer {
		return nil, errs.NewInternalServerError("cursor res param is not pointer")
	}

	db := r.DB(ctx).Table(tableName)
	db = parseFilter(db, req)

	// 获取游标参数
	cursor := 0
	count := 10 // 默认每页大小
	if c, ok := req.(search.Cursor); ok {
		cursor = c.GetCursor()
		count = c.GetCount()
	}

	err := db.
		Offset(cursor).
		Limit(count + 1). // 多查询一条用于判断是否还有更多数据
		Find(res).
		Error

	if err != nil {
		return nil, ErrorMap(err, tableName)
	}

	valueOf = reflect.ValueOf(res)
	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}

	if valueOf.Kind() != reflect.Slice {
		return nil, errs.NewInternalServerError("cursor res must be slice")
	}

	resultCount := valueOf.Len()
	hasMore := resultCount > count
	nextCursor := cursor

	// 如果查询到的数据超过请求的数量，说明还有更多数据
	if hasMore {
		// 删除多查询的那条数据
		valueOf.Set(valueOf.Slice(0, count))
		nextCursor += count
	} else {
		nextCursor += resultCount
	}

	return search.NewCursorResult(nextCursor, hasMore), nil
}

// Cursor 游标分页查询
func (r *DAOImpl) CursorByCondition(ctx context.Context, tableName string, req search.Cursor, res any) (*search.CursorResult, error) {
	valueOf := reflect.ValueOf(res)
	if valueOf.Kind() != reflect.Pointer {
		return nil, errs.NewInternalServerError("cursor res param is not pointer")
	}

	condition := search.MakeCursorCondition(tableName, req.GetCursor(), req.GetCount(), req)

	// 查询数据
	db := r.DB(ctx)
	db = r.applyCondition(db, condition, false)
	err := db.
		Table(tableName).
		Offset(req.GetCursor()).
		Limit(req.GetCount() + 1). // 多查询一条用于判断是否还有更多数据
		Find(res).
		Error

	if err != nil {
		return nil, ErrorMap(err, tableName)
	}

	valueOf = reflect.ValueOf(res)
	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}

	if valueOf.Kind() != reflect.Slice {
		return nil, errs.NewInternalServerError("cursor res must be slice")
	}

	count := valueOf.Len()
	hasMore := count > req.GetCount()
	nextCursor := req.GetCursor()

	// 如果查询到的数据超过请求的数量，说明还有更多数据
	if hasMore {
		// 删除多查询的那条数据
		valueOf.Set(valueOf.Slice(0, req.GetCount()))
		nextCursor += req.GetCount()
	} else {
		nextCursor += count
	}

	return search.NewCursorResult(nextCursor, hasMore), nil
}

func (r *DAOImpl) applyCondition(db *gorm.DB, condition search.Condition, enableCount bool) *gorm.DB {
	rootDB := db
	projection := condition.GetProjection()
	if projection != nil && len(projection) > 0 {
		db = db.Select(projection)
	}

	distinct := condition.GetDistinct()
	if distinct != nil && len(distinct) > 0 {
		db = db.Distinct(toInterfaceSlice(distinct)...)
	}

	dbCondition := NewDBCondition(condition)

	db = r.parseCondition(rootDB, db, dbCondition, enableCount)

	return db
}

func (r *DAOImpl) parseCondition(rootDB *gorm.DB, db *gorm.DB, dbCondition *GormCondition, enableCount bool) *gorm.DB {
	// group or
	if dbCondition.GroupOr != nil {
		db = db.Where(r.parseCondition(rootDB, rootDB, dbCondition.GroupOr, enableCount))
	}

	for k, v := range dbCondition.Where {
		db = db.Where(k, v...)
	}

	for k, v := range dbCondition.Or {
		db = db.Or(k, v...)
	}

	for _, join := range dbCondition.Join {
		if join == nil {
			continue
		}
		query := join.On
		var args []any

		for k, v := range join.Where {
			query = fmt.Sprintf("%s AND %s", query, k)
			args = append(args, v...)
		}
		for k, v := range join.Or {
			query = fmt.Sprintf("%s OR %s", query, k)
			args = append(args, v...)
		}

		db = db.Joins(query, args...)
	}

	if !enableCount {
		for _, o := range dbCondition.Order {
			db = db.Order(o)
		}
	}

	for _, o := range dbCondition.Group {
		db = db.Group(o)
	}

	return db
}

func toInterfaceSlice(strings []string) []any {
	interfaces := make([]any, len(strings))
	for i, s := range strings {
		interfaces[i] = s
	}
	return interfaces
}

// FindInBatchesV2 完全避开 GORM 的 FindInBatches，手动实现分批查询
func (r *DAOImpl) FindInBatchesV2(ctx context.Context, tableName string, req any, res any, batchSize int, fn func(batch int) error) error {
	// 参数验证
	if batchSize <= 0 {
		return fmt.Errorf("batchSize must be positive, got %d", batchSize)
	}
	if fn == nil {
		return fmt.Errorf("callback function cannot be nil")
	}

	// 检查是否为 interfaces.Filter 类型（动态查询）
	if filterReq, ok := req.(interfaces.Filter); ok {
		// 使用 parseFilter 处理动态查询，类似 FindInBatches
		return r.findInBatchesWithFilter(ctx, tableName, filterReq, res, batchSize, fn)
	}

	// 验证 res 参数并获取反射值（移到循环外部以提高性能）
	resValue := reflect.ValueOf(res)
	if resValue.Kind() != reflect.Ptr {
		return fmt.Errorf("res must be a pointer, got %T", res)
	}

	resElem := resValue.Elem()
	if resElem.Kind() != reflect.Slice {
		return fmt.Errorf("res must be a pointer to slice, got pointer to %s", resElem.Kind())
	}

	// 获取切片类型和初始容量
	sliceType := resElem.Type()
	initialCap := resElem.Cap()
	if initialCap == 0 {
		initialCap = batchSize
	}

	// 使用结构化查询
	condition := search.MakeCondition(tableName, req)
	offset := 0
	batch := 0

	for {
		// 安全地清空切片，保持容量以避免频繁内存分配
		resElem.Set(reflect.MakeSlice(sliceType, 0, initialCap))

		db := r.DB(ctx).Table(tableName)
		db = r.applyCondition(db, condition, false)
		db = db.Limit(batchSize).Offset(offset)

		result := db.Find(res)
		if result.Error != nil {
			return ErrorMap(result.Error, tableName)
		}

		rowsAffected := int(result.RowsAffected)
		if rowsAffected == 0 {
			break
		}

		// 执行回调函数，传入当前批次号（从0开始）
		if err := fn(batch); err != nil {
			return fmt.Errorf("callback function failed at batch %d: %w", batch, err)
		}

		// 如果返回的记录数少于批次大小，说明已经是最后一批
		if rowsAffected < batchSize {
			break
		}

		offset += batchSize
		batch++
	}

	return nil
}

// findInBatchesWithFilter 处理 interfaces.Filter 类型的批量查询
func (r *DAOImpl) findInBatchesWithFilter(ctx context.Context, tableName string, req interfaces.Filter, res any, batchSize int, fn func(batch int) error) error {
	offset := 0
	batch := 0

	for {
		db := r.DB(ctx).Table(tableName)
		db = parseFilter(db, req)
		db = db.Limit(batchSize).Offset(offset)

		result := db.Find(res)
		if result.Error != nil {
			return ErrorMap(result.Error, tableName)
		}

		rowsAffected := int(result.RowsAffected)
		if rowsAffected == 0 {
			break
		}

		// 执行回调函数，传入当前批次号（从0开始）
		if err := fn(batch); err != nil {
			return fmt.Errorf("callback function failed at batch %d: %w", batch, err)
		}

		// 如果返回的记录数少于批次大小，说明已经是最后一批
		if rowsAffected < batchSize {
			break
		}

		offset += batchSize
		batch++
	}

	return nil
}
