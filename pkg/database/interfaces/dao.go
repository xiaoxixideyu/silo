package interfaces

import (
	"context"
	"silo/pkg/database/search"
	"silo/pkg/types/isotime"

	"gorm.io/gorm"
)

//go:generate mockgen -source=dao.go -destination=mocks/dao.go

// DaoRead only
type DaoRead interface {
	// Query
	FindOne(ctx context.Context, tableName string, req Filter, res any) error
	FindOneByCondition(ctx context.Context, tableName string, req any, res any) error

	Find(ctx context.Context, tableName string, req Filter, res any) error
	FindByCondition(ctx context.Context, tableName string, req any, res any) error

	Page(ctx context.Context, tableName string, req search.Pagination, res any) (*search.PageResult, error)

	FindInBatches(ctx context.Context, tableName string, req Filter, res any, batchSize int, fn func(batch int) error) error
	FindInBatchesByCondition(ctx context.Context, tableName string, req any, res any, batchSize int, fn func(batch int) error) error
	FindInBatchesV2(ctx context.Context, tableName string, req any, res any, batchSize int, fn func(batch int) error) error

	// Count
	Count(ctx context.Context, tableName string, req Filter) (int64, error)
	CountByCondition(ctx context.Context, tableName string, req any) (int64, error)

	Cursor(ctx context.Context, tableName string, req Filter, res any) (*search.CursorResult, error)
	CursorByCondition(ctx context.Context, tableName string, req search.Cursor, res any) (*search.CursorResult, error)
}

// Dao .
type Dao interface {
	DaoRead
	// Create
	Create(ctx context.Context, tableName string, model Model) error
	CreateMany(ctx context.Context, tableName string, models any) error
	CreateInBatches(ctx context.Context, tableName string, models any, batchSize int) error
	FindOrCreate(ctx context.Context, tableName string, defaults Model, res any) (create bool, err error)

	// Update
	UpdateMany(ctx context.Context, model Model, req Update) error
	FindOneAndUpdate(ctx context.Context, tableName string, req Update, res any) error

	// Save
	Save(ctx context.Context, model Model) error
	SaveInBatches(ctx context.Context, models []Model) error

	// DeleteOne
	Delete(ctx context.Context, tableName string, model Model) error
	DeleteOne(ctx context.Context, tableName string, req Filter) error
	DeleteMany(ctx context.Context, tableName string, req Filter) error

	// IncreaseAndReturn
	IncreaseAndReturn(ctx context.Context, tableName string, req Filter, column string, step int64, res any) error
}

// PostgresDao .
type PostgresDao interface {
	Dao
	DB(ctx context.Context) *gorm.DB
}

// PostgresDAORead .
type PostgresDAORead interface {
	Dao
	DB(ctx context.Context) *gorm.DB
}

/************ With ************/

type WithTimestamp interface {
	GetUpdatedAt() isotime.ISOTime
	WithTimestampSetUpdatedAt()
	GetCreatedAt() isotime.ISOTime
	WithTimestampSetCreatedAt()
}

type WithSoftDelete interface {
	SetDeletedAt()
}

/************ Request ************/

const ASC = 1
const DESC = -1

type M = map[string]any
type D = []E
type E = struct {
	Key   string
	Value any
}
type A = []any

type Filter interface {
	GetFilter() D
	GetLimit() int64
	GetSort() D
	GetProject() D
	GetGroup() []string
}

type Upsert interface {
	GetUpsert() bool
}

type Update interface {
	Filter
	Upsert
	GetUpdate() M
}
