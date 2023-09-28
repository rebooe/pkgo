package cache

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

type Cacher interface {
	Get(string, any) (bool, error)
	Set(string, any, time.Duration) error
	Delete(string) error
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

func setValue(tag, src any) error {
	// 使用反射设置 tag 的值
	tagValue := reflect.ValueOf(tag)
	// 确保 tag 是一个指针
	if tagValue.Kind() != reflect.Ptr {
		return errors.New("tag must be a pointer")
	}
	tagValue = tagValue.Elem()
	// 获取 tag 的类型
	tagType := tagValue.Type()

	// 将 src 转换为 tag 的类型并设置给 tag
	if !reflect.TypeOf(src).ConvertibleTo(tagType) {
		return errors.New("src cannot be converted to the type of tag")
	}
	tagValue.Set(reflect.ValueOf(src).Convert(tagType))
	return nil
}
