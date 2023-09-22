package logger

import (
	"log"
	"os"
	"testing"
)

func TestLog_Printf(t *testing.T) {
	w, err := NewRouterWriter("", "test")
	if err != nil {
		t.Fatal(err)
	}

	logger := log.New(w, "[TEST]", log.LstdFlags)
	logger.Print("AAA")
}

func TestLog_Delete(t *testing.T) {
	os.RemoveAll("test")
}
