package ioc

import (
	"github.com/redis/go-redis/v9"
	"github/yyfzy/mybook/config"
)

func InitRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return cmd
}
