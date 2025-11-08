package scripts

import (
	"context"
	"errors"
	"silo/pkg/database/interfaces"
	"silo/pkg/errs"
	"strings"

	"github.com/redis/go-redis/v9"
)

// CacheNotFoundError 缓存不存在错误
const CacheNotFoundError = redis.Nil

// 错误码
const (
	ErrOptimisticLockError = "800" // 加锁中异常
)

// ErrorMapper Redis错误转换函数
func ErrorMapper(cacheClient interfaces.CacheClient, err error) error {
	if err == nil {
		return nil
	}

	if cacheClient.IsNil(err) {
		return errs.NewInstanceNotFoundError(err.Error())
	}

	if errors.Is(err, context.Canceled) {
		return err
	}

	msg := err.Error()
	msg = strings.TrimLeft(msg, "ERR ")

	switch msg {
	case ErrOptimisticLockError:
		return errs.NewCustomError(errs.StatusOptimisticLockError, "Optimistic lock failed")
	}

	return errs.NewInternalServerErrorf("cache error: %s", msg)
}
