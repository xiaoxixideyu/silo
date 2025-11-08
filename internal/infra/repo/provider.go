package repo

import (
	"silo/pkg/database/dao"
	"silo/pkg/database/db_impl/postgres"
	"silo/pkg/database/db_impl/redis"
	"silo/pkg/database/interfaces"
	"silo/pkg/database/transaction"

	"github.com/google/wire"
)

var RepoSet = wire.NewSet(
	postgres.NewPostgres,
	redis.NewRedisClient,
	dao.NewDao,
	dao.NewDaoCache,
	dao.NewDaoRead,
	dao.NewTXBeginner,
	transaction.NewManager,
	NewExampleRepo,
	wire.Bind(new(interfaces.CacheClient), new(*redis.RedisClient)),
)
