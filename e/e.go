package e

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Warp 包装堆栈信息到错误
func Warp(err error) error {
	if err == nil {
		return nil
	}
	// 获取堆栈信息
	pc, file, line, _ := runtime.Caller(1)
	baseFile := filepath.Base(file)
	Func := runtime.FuncForPC(pc)
	pack, funcName, _ := strings.Cut(Func.Name(), ".")

	return fmt.Errorf("%s/%s:%d %s()\n%w", pack, baseFile, line, funcName, err)
}

// Warpf 包装堆栈信息到格式化错误
func Warpf(format string, args ...any) error {
	return Warp(fmt.Errorf(format, args...))
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
