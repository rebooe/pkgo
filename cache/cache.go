package cache

import (
	"time"
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

