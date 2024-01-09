package cache

import (
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

type goCache struct {
	*cache.Cache
}

func NewGoCache() Cacher {
	return &goCache{
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (c *goCache) Get(key string) Valuer {
	value, expir, exist := c.GetWithExpiration(key)
	return &gocacheValue{
		v:     value,
		expir: expir,
		exist: exist,
	}
}

func (c *goCache) Set(key string, val any, duration time.Duration) error {
	c.Cache.Set(key, val, duration)
	return nil
}

func (c *goCache) Delete(key string) error {
	c.Cache.Delete(key)
	return nil
}

type gocacheValue struct {
	v     any
	expir time.Time
	exist bool
}

func (v *gocacheValue) Any() (any, error) { return v.v, v.Err() }

func (v *gocacheValue) Int() (int, error) {
	if !v.exist {
		return 0, nil
	}

	to, ok := v.v.(int)
	if !ok {
		return 0, errors.New("type is not int")
	}
	return to, nil
}

func (v *gocacheValue) Float() (float64, error) {
	if !v.exist {
		return 0.0, nil
	}

	to, ok := v.v.(float64)
	if !ok {
		return 0.0, errors.New("type is not float64")
	}
	return to, nil
}

func (v *gocacheValue) String() (string, error) {
	if !v.exist {
		return "", nil
	}

	to, ok := v.v.(string)
	if !ok {
		return "", errors.New("type is not string")
	}
	return to, nil
}

func (v *gocacheValue) Byte() ([]byte, error) {
	if !v.exist {
		return nil, nil
	}

	to, ok := v.v.([]byte)
	if !ok {
		return nil, errors.New("type is not []byte")
	}
	return to, nil
}

func (v *gocacheValue) Err() error { return nil }

func (v *gocacheValue) Expir() time.Time { return v.expir }
