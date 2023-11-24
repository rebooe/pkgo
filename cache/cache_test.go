package cache

import (
	"fmt"
	"log"
	"testing"
	"time"
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

	value, err := c.Get("test").ToAny()
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
