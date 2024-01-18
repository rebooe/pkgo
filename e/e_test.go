package e

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := Warp(errors.New("err"), "123")
	t.Logf("%s", err)

	err = Warp(err)
	t.Logf("%s", err)
}
