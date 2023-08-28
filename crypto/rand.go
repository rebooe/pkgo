package crypto

import (
	"crypto/rand"
	"io"
)

func RandByte(length int) ([]byte, error) {
	r := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, r); err != nil {
		return nil, err
	}
	return r, nil
}
