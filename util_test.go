package pkgo

import (
	"testing"
)

func TestCaller(t *testing.T) {
	got, _ := Caller(0)
	t.Logf("skip:0, %s", got)

	func() {
		got, _ := Caller(0)
		t.Logf("skip:0, %s", got)
	}()
}
