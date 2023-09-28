package util

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"
)

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Caller 获取格式化的一条堆栈信息
func Caller(skip int) (string, bool) {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "", false
	}
	baseFile := filepath.Base(file)
	Func := runtime.FuncForPC(pc)
	pack, funcName, _ := strings.Cut(Func.Name(), ".")

	return fmt.Sprintf("%s/%s:%d %s()", pack, baseFile, line, funcName), true
}

// Callers 获取格式化的所有堆栈信息
func Callers() (strs []string) {
	pc := make([]uintptr, 32)
	runtime.Callers(3, pc)

	frames := runtime.CallersFrames(pc)
	for {
		if frame, more := frames.Next(); more {
			file := filepath.Base(frame.File)
			pack, fn, _ := strings.Cut(frame.Function, ".")

			strs = append(strs, fmt.Sprintf("%s/%s:%d %s()", pack, file, frame.Line, fn))
			continue
		}
		break
	}
	return
}

// ReplacePlaceholders 替换查询字符串中的 ? 为实际参数值
func ReplacePlaceholders(query string, args ...interface{}) string {
	query = strings.ReplaceAll(query, "?", "%v")
	return fmt.Sprintf(query, args...)
}
