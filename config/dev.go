//go:build !k8s

package config

var Config = config{
	DB:    DBConfig{DSN: "root:root@tcp(192.168.137.132:13316)/webook"},
	Redis: RedisConfig{Addr: "192.168.137.132:6379"},
}
