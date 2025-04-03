package cache

import (
	"context"
	"testing"
)

func TestNewCache(t *testing.T) {
	cache, err := NewCache(&CacheConfig{
		Type: "store",
	})
	if err != nil {
		t.Error(err)
	}

	if err := cache.Set(context.Background(), "test", "test"); err != nil {
		t.Error(err)
	}

	value, err := cache.Get(context.Background(), "test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", value)
}
