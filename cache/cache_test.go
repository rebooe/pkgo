package cache

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type test struct {
	A int
	B string
}

func testCache(c Cacher) error {
	err := c.Set("test", &test{A: 123, B: "123"}, time.Hour)
	if err != nil {
		return err
	}

	value, err := c.Get("test").Any()
	if err != nil {
		return err
	}

	res, _ := value.(*test)
	if res.A != 123 {
		return fmt.Errorf("A got %v, want %v", res.A, 123)
	}
	if res.B != "123" {
		return fmt.Errorf("B got %v, want %v", res.B, "123")
	}

	log.Printf("res: %v", res)
	return nil
}

func Test_goCache(t *testing.T) {
	c := NewGoCache()
	if err := testCache(c); err != nil {
		t.Fatal(err)
	}
}

func Test_redisCache(t *testing.T) {
	c := NewRedisCache(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	if err := testCache(c); err != nil {
		t.Fatal(err)
	}
}
