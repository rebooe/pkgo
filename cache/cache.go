package cache

import (
	"errors"
	"time"

	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	gocache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

type CacheConfig struct {
	Type      string `yaml:"Type"`      // 缓存类型(store:本地缓存,redis:redis缓存)
	RedisAddr string `yaml:"RedisAddr"` // redis地址
	RedisPwd  string `yaml:"RedisPwd"`  // redis密码
	RedisDB   int    `yaml:"RedisDB"`   // redis数据库
}

func NewCache(config *CacheConfig) (store.StoreInterface, error) {
	switch config.Type {
	case "store":
		gocacheClient := gocache.New(5*time.Minute, 10*time.Minute)
		return gocache_store.NewGoCache(gocacheClient), nil
	case "redis":
		return redis_store.NewRedis(
			redis.NewClient(&redis.Options{
				Addr:     config.RedisAddr,
				Password: config.RedisPwd,
				DB:       config.RedisDB,
			})), nil
	}
	return nil, errors.New("cache type not found")
}
