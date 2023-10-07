package e

import (
	"testing"
)

func TestNew(t *testing.T) {
	err := Warpf("123")
	t.Logf("%s", err)

	err = Warp(err)
	t.Logf("%s", err)
}
