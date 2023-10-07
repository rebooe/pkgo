package e

import (
	"fmt"

	"github.com/rebooe/pkg-go/util"
)

// Warp 包装堆栈信息到错误
func Warp(err error) error {
	if err == nil {
		return nil
	}
	s, _ := util.Caller(1)
	return fmt.Errorf("%s\n%w", s, err)
}

// Warpf 包装堆栈信息到格式化错误
func Warpf(format string, args ...any) error {
	s, _ := util.Caller(1)
	return fmt.Errorf("%s\n%w", s, fmt.Errorf(format, args...))
}

// Cause 返回根错误
func Cause(err error) error {
	for {
		e, ok := err.(interface{ Unwrap() error })
		if !ok {
			return err
		}
		err = e.Unwrap()
		if err == nil {
			return nil
		}
	}
}
