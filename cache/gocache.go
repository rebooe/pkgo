package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type goCache struct {
	*cache.Cache
}

var _ Cacher = (*goCache)(nil)

func NewGoCache() Cacher {
	return &goCache{
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (c *goCache) Get(key string, v any) (bool, time.Time, error) {
	value, expir, exist := c.GetWithExpiration(key)
	if !exist {
		return false, expir, nil
	}
	err := setValue(v, value)
	return true, expir, err
}

func (c *goCache) Set(key string, val any, duration time.Duration) error {
	c.Cache.Set(key, val, duration)
	return nil
}

func (c *goCache) Delete(key string) error {
	c.Cache.Delete(key)
	return nil
}
