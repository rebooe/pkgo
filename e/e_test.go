package e

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := Wrap(errors.New("err"), "123")
	t.Logf("%s", err)

	err = Wrap(err)
	t.Logf("%s", err)
}
