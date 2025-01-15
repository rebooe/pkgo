package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(opts *redis.Options) Cacher {
	client := redis.NewClient(opts)
	return &redisCache{client: client}
}

func (c *redisCache) Get(key string) Valuer {
	cmd := c.client.Get(context.Background(), key)
	ttl := c.client.TTL(context.Background(), key).Val()
	return &redisValue{
		cmd:   cmd,
		expir: time.Now().Add(ttl),
	}
}

func (c *redisCache) Set(key string, value any, t time.Duration) error {
	return c.client.Set(context.Background(), key, value, t).Err()
}

func (c *redisCache) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

type redisValue struct {
	cmd   *redis.StringCmd
	expir time.Time
}

func (v *redisValue) Any() (any, error) {
	return v.cmd.Result()
}

func (v *redisValue) Int() (int, error) {
	return v.cmd.Int()
}

func (v *redisValue) Float() (float64, error) {
	return v.cmd.Float64()
}

func (v *redisValue) String() (string, error) {
	return v.cmd.Result()
}

func (v *redisValue) Byte() ([]byte, error) {
	return v.cmd.Bytes()
}

func (v *redisValue) Err() error {
	return v.cmd.Err()
}

func (v *redisValue) Expir() time.Time {
	return v.expir
}

func (v *redisValue) Exists() bool {
	return v.cmd.Val() != ""
}
