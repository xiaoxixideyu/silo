package scripts

import (
	// for import lua
	"context"
	_ "embed"
	"silo/pkg/database/interfaces"

	"github.com/redis/go-redis/v9"
)

//go:embed id_generate_get.lua
var luaIDGenerateGet string
var idGenerateGet = redis.NewScript(luaIDGenerateGet)

func IDGenerateGet(ctx context.Context, cacheClient interfaces.CacheClient, key string, keepAliveKey string, maxID int, timeout int) (int, error) {
	id, err := idGenerateGet.Run(ctx, cacheClient, []string{key, keepAliveKey}, maxID, timeout).Int()
	if err != nil {
		return 0, err
	}
	return id, nil
}

//go:embed id_generate_free.lua
var luaIDGenerateFree string
var idGenerateFree = redis.NewScript(luaIDGenerateFree)

func IDGenerateFree(ctx context.Context, cacheClient interfaces.CacheClient, key string, keepAliveKey string, id int) error {
	_, err := idGenerateFree.Run(ctx, cacheClient, []string{key, keepAliveKey}, id).Int()
	if err != nil {
		return err
	}
	return nil
}

//go:embed id_generate_keepalive.lua
var luaIDGenerateKeepalive string
var idGenerateKeepalive = redis.NewScript(luaIDGenerateKeepalive)

func IDGenerateKeepalive(ctx context.Context, cacheClient interfaces.CacheClient, keepAliveKey string, id int) error {
	_, err := idGenerateKeepalive.Run(ctx, cacheClient, []string{keepAliveKey}, id).Int()
	if err != nil {
		return err
	}
	return nil
}
