package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"silo/pkg/database/config"
	"strings"
	"time"

	apmgoredis "github.com/ekrucio/apm-agent-go/module/apmgoredisv9/v2"
	"github.com/redis/go-redis/v9"
)

// RedisClient Client
type RedisClient struct {
	redis.UniversalClient
	prefix string
}

const CacheNotFoundError = redis.Nil

// NewRedisClient client
func NewRedisClient(cfg config.RedisConfig) *RedisClient {
	var client *RedisClient
	var rc redis.UniversalClient

	if cfg.ClusterEnabled {
		nodes := strings.Split(cfg.Addr, ",")
		rc = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    nodes,
			Password: cfg.Password,
		})

	} else {
		rc = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password, // no password set
			DB:       0,            // use default DB
		})
	}

	client = &RedisClient{
		UniversalClient: rc,
		prefix:          cfg.Prefix,
	}

	rc.AddHook(apmgoredis.NewHook())

	return client
}

// IsNil check if the error is a cache not found error
func (c *RedisClient) IsNil(err error) bool {
	return errors.Is(err, redis.Nil)
}

// Key format cache key with prefix
func (c *RedisClient) Key(key string) string {
	return c.prefix + key
}

// Keyf format cache key with prefix and args
func (c *RedisClient) Keyf(keyTemplate string, args ...any) string {
	k := c.prefix + keyTemplate
	return fmt.Sprintf(k, args...)
}

// UnmarshalGet UnmarshalGet
func (c *RedisClient) UnmarshalGet(ctx context.Context, key string, out any) (bool, error) {
	str, err := c.Get(ctx, key).Result()
	if err != nil {
		if c.IsNil(err) {
			return false, nil
		}
		return false, err
	}

	e := json.Unmarshal([]byte(str), out)
	if e != nil {
		return false, e
	}

	return true, nil
}

// MarshalSet MarshalSet
func (c *RedisClient) MarshalSet(ctx context.Context, key string, value any, expiration time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, e := c.Set(ctx, key, v, expiration).Result()
	if e != nil {
		return e
	}

	return nil
}

// UnmarshalHGet UnmarshalHGet
func (c *RedisClient) UnmarshalHGet(ctx context.Context, key string, field string, out any) (bool, error) {
	str, err := c.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}

	e := json.Unmarshal([]byte(str), out)
	if e != nil {
		return false, e
	}

	return true, nil
}

// MarshalHSet .
func (c *RedisClient) MarshalHSet(ctx context.Context, key string, field string, value any) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, e := c.HSet(ctx, key, field, v).Result()
	if e != nil {
		return e
	}

	return nil
}

// ZAddHandle .
func (c *RedisClient) ZAddHandle(ctx context.Context, key string, score float64, member any) (int64, error) {
	return c.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Result()
}
