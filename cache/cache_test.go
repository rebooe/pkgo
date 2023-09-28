package cache

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

type test struct {
	A int
	B string
}

func testCache(c Cacher) error {
	err := c.Set("test", test{A: 123, B: "123"}, time.Hour)
	if err != nil {
		return err
	}

	var res test
	ok, err := c.Get("test", &res)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("ok is %v", ok)
	}
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

func Test_SqliteCache(t *testing.T) {
	db, err := sql.Open("sqlite", "cache.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewSqliteCacher(db, "cache")
	if err != nil {
		t.Fatal(err)
	}

	if err := testCache(c); err != nil {
		t.Fatal(err)
	}
	os.Remove("cache.sqlite")
}
