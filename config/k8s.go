//go:build k8s

package config

var Config = config{
	DB:    DBConfig{DSN: "webook-mysql:13309"},
	Redis: RedisConfig{Addr: "webook-redis:11479"},
}
