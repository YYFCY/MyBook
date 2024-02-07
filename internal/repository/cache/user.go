package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github/yyfzy/mybook/internal/domain"
	"time"
)

var ErrKeyNotExit = redis.Nil

type CacheV1 interface {
	GetUser(ctx context.Context, id int64) domain.User
}

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id int64) error
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

// A用到B, B一定是接口
// A用到B, B一定是A的字段
// A用到B, A绝对不初始化B, 而是从外面注入
func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.Key(id)
	val, err := cache.client.Get(ctx, key).Result()
	fmt.Println(key)
	fmt.Println(val, err)
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(val), &u)
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.Key(u.Id)
	fmt.Println(cache.client)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()

}

func (cache *RedisUserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (cache *RedisUserCache) Delete(ctx context.Context, id int64) error {
	return cache.client.Del(ctx, cache.Key(id)).Err()
}

//type UnifyCache interface {
//	Get(ctx context.Context, key string)
//	Set(ctx context.Context, key string, val any, expiration time.Duration)
//}
