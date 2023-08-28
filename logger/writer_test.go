package logger

import (
	"log"
	"testing"
)

func TestLog_Printf(t *testing.T) {
	w, err := New("", "test")
	if err != nil {
		t.Fatal(err)
	}

	logger := log.New(w, "[TEST]", log.LstdFlags)
	logger.Print("AAA")
}
