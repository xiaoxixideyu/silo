package dao

import (
	"context"
	"errors"
	"fmt"
	"silo/pkg/database/interfaces"
	"silo/pkg/errs"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type dbKey struct{}

// DBToCtx .
func DBToCtx(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey{}, db)
}

// CtxToDB .
func CtxToDB(ctx context.Context) *gorm.DB {
	val := ctx.Value(dbKey{})
	if val == nil {
		return nil
	}

	db, ok := val.(*gorm.DB)
	if !ok {
		return nil
	}

	return db
}

// NewDao .
func NewDao(db *gorm.DB) interfaces.Dao {
	SetEnclose("\"")
	return &DAOImpl{
		db: db,
	}
}

// NewDaoRead .
func NewDaoRead(db *gorm.DB) interfaces.DaoRead {
	SetEnclose("\"")
	return &DAOImpl{
		db: db,
	}
}

// NewTXBeginner .
func NewTXBeginner(db *gorm.DB) interfaces.TXBeginner {
	return &DAOImpl{
		db: db,
	}
}

// DAOImpl .
type DAOImpl struct {
	db *gorm.DB
}

// ErrorMap .
func ErrorMap(err error, tableName string) error {
	if err != nil {
		if errs.IsCustomError(err) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewInstanceNotFoundError(err.Error()).WithError(fmt.Errorf("table %s", tableName))
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.NewDuplicatedError(err.Error()).WithError(fmt.Errorf("table %s", tableName))
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return errs.NewDuplicatedError(err.Error()).WithError(fmt.Errorf("table %s", tableName))
			}
		} else if errors.Is(err, context.DeadlineExceeded) {
			return errs.NewRequestTimeoutError(err.Error()).WithError(fmt.Errorf("table %s", tableName))
		} else if errors.Is(err, context.Canceled) {
			return errs.NewRequestCanceledError(err.Error()).WithError(fmt.Errorf("table %s", tableName))
		}
		return errs.NewInternalServerError(err.Error()).WithError(fmt.Errorf("table %s, error %w", tableName, err))
	}

	return nil
}
