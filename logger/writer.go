package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type RotateWrite struct {
	writer io.WriteCloser
	dir    string // 日志存放路径
	date   string // 当前日期
}

func New(dir, tagName string) (io.Writer, error) {
	dir = filepath.Join(dir, tagName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &RotateWrite{dir: dir}, nil
}

func (w *RotateWrite) Write(p []byte) (int, error) {
	now := time.Now().Format("20060102")

	if now != w.date {
		// 日志滚动
		file := filepath.Join(w.dir, now+".log")
		f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return 0, fmt.Errorf("创建新日志文件错误: %w", err)
		}
		if w.writer != nil {
			if err := w.writer.Close(); err != nil {
				return 0, fmt.Errorf("关闭旧日志文件错误: %w", err)
			}
		}

		w.writer = f
		w.date = now
	}

	return w.writer.Write(p)
}
