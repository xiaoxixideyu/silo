package interfaces

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=cache.go -destination=mocks/cache.go

// CacheKeyFormat .
type CacheKeyFormat interface {
	// apply redis key prefix
	Key(key string) string
	// apply redis key prefix to template
	Keyf(keyTemplate string, args ...any) string
}

// CacheClient .
type CacheClient interface {
	redis.UniversalClient
	CacheKeyFormat
	UnmarshalGet(ctx context.Context, key string, out any) (bool, error)
	MarshalSet(ctx context.Context, key string, value any, expiration time.Duration) error
	UnmarshalHGet(ctx context.Context, key string, field string, out any) (bool, error)
	MarshalHSet(ctx context.Context, key string, field string, value any) error
	ZAddHandle(ctx context.Context, key string, score float64, member any) (int64, error)
	IsNil(err error) bool
}

// Pipeliner .
type Pipeliner interface {
	redis.Pipeliner
}

type RedisPubsubOptions struct {
	DisableLog bool
}

type RedisPubsubOption func(*RedisPubsubOptions)

func RedisPubsubWithDisableLog() RedisPubsubOption {
	return func(o *RedisPubsubOptions) {
		o.DisableLog = true
	}
}

// RedisPubsub .
type RedisPubsub interface {
	CacheKeyFormat
	Publish(ctx context.Context, channel string, message string) error
	Subscribe(ctx context.Context, channel string, handler func(ctx context.Context, channel, message string) error, options ...RedisPubsubOption)
}
