package redis

import (
	"bluebell/settings"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	ctx context.Context
	rdb *redis.Client
)

func Init(rcf *settings.RedisConfig) (err error) {
	// rdb = redis.NewClient(&redis.Options{
	// 	Addr: fmt.Sprintf("%s:%d",
	// 		viper.GetString("redis.host"),
	// 		viper.GetInt("redis.port"),
	// 	),
	// 	Password: viper.GetString("redis.password"),
	// 	DB:       viper.GetInt("redis.db"),
	// 	PoolSize: viper.GetInt("redis.pool_size"),
	// })
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			rcf.Host,
			rcf.Port,
		),
		Password: rcf.Password,
		DB:       rcf.Db,
		PoolSize: rcf.PoolSize,
	})
	ctx = rdb.Context()
	_, err = rdb.Ping(ctx).Result()
	return
}
func Close() {
	rdb.Close()
}
