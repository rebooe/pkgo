package queue

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"time"
)

type FileMQ struct {
	dir        string
	err        chan error
	pushSignal chan struct{} // 任务推送信号
}

func NewFileMQ(dir string) (Queue, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &FileMQ{
		dir:        dir,
		err:        make(chan error),
		pushSignal: make(chan struct{}),
	}, nil
}

func (mq *FileMQ) Push(name string, msg Msg) error {
	if err := os.MkdirAll(path.Join(mq.dir, name), 0755); err != nil {
		return err
	}
	ID := fmt.Sprintf("%x", time.Now().UnixNano())

	buf := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(msg); err != nil {
		return err
	}

	fileName := path.Join(mq.dir, name, ID)
	if err := os.WriteFile(fileName, buf.Bytes(), 0655); err != nil {
		return err
	}

	// 触发推送任务信号
	select {
	case mq.pushSignal <- struct{}{}:
	default:
	}
	return nil
}

func (mq *FileMQ) Pop(name string, opts ...OptionFunc) (*Msg, error) {
	opt := defaultOption()
	for i := range opts {
		opts[i](opt)
	}

	// 读取目录下任务
	dirEntry, err := os.ReadDir(path.Join(mq.dir, name))
	if err != nil {
		return nil, err
	}

	if len(dirEntry) == 0 {
		// 无任务时判断是否阻塞
		if !opt.Wait {
			return nil, nil
		}
		<-mq.pushSignal
	}

	// 读取目录下任务
	dirEntry, err = os.ReadDir(path.Join(mq.dir, name))
	if err != nil {
		return nil, err
	}

	file := dirEntry[0]
	body, err := os.ReadFile(path.Join(mq.dir, name, file.Name()))
	if err != nil {
		return nil, err
	}

	var msg Msg
	decoder := gob.NewDecoder(bytes.NewBuffer(body))
	if err := decoder.Decode(&msg); err != nil {
		return nil, err
	}

	// 是否确认收到
	if opt.Ack {
		if err := os.Remove(path.Join(mq.dir, name, file.Name())); err != nil {
			return nil, err
		}
	}
	return &msg, nil
}
