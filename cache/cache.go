package cache

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type Cacher interface {
	Get(string) Valuer
	Set(string, any, time.Duration) error
	Delete(string) error
}

type Valuer interface {
	Any() (any, error)
	Int() (int, error)
	Float() (float64, error)
	String() (string, error)
	Byte() ([]byte, error)

	Err() error
	Expir() time.Time
}

const cacheKey = "cache/cacheKey"

// 设置缓存到上下文
func SetCache(cacher Cacher) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(cacheKey, cacher)
	}
}

// 从上下文中获取缓存
func GetCache(c context.Context) Cacher {
	return c.Value(cacheKey).(Cacher)
}
