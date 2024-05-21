package queue

import (
	"testing"
)

var q Queue

func TestMain(m *testing.M) {
	fileMQ, err := NewFileMQ("test")
	if err != nil {
		panic(err)
	}
	q = fileMQ
	m.Run()
}

func TestFileMQ_Push(t *testing.T) {
	for i := 0; i < 10; i++ {
		err := q.Push("", Msg{Body: i})
		if err != nil {
			panic(err)
		}
	}
}

func TestFileMQ_Pop(t *testing.T) {
	msg, err := q.Pop("", OptAck(true), OptWait(false))
	if err != nil {
		panic(err)
	}

	t.Logf("%v", msg)
}
