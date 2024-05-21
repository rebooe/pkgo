package e

import (
	"fmt"
	"strings"

	"github.com/rebooe/pkgo"
)

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

type myErr struct {
	callers []string // 堆栈信息
	err     error
}

func (e *myErr) Error() string {
	return fmt.Sprintf("%s\n%s", strings.Join(e.callers, "\n"), e.err)
}

func (e *myErr) Unwrap() error {
	return e.err
}

// Warp 包装错误
//
//	err 原始错误
//	msg 额外信息
func Warp(err error, msg ...string) error {
	if err == nil {
		return nil
	}

	s, _ := pkgo.Caller(1)
	e, ok := err.(*myErr)
	if ok {
		e.callers = append(e.callers, s)
	} else {
		e = &myErr{
			callers: []string{s},
			err:     err,
		}
	}

	if len(msg) > 0 {
		e.err = fmt.Errorf("%s: %s", msg[0], e.err)
	}
	return e
}
