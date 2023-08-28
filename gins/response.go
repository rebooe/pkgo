package gins

import (
	"net/http"
	"pkg-go/e"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code  int
	Msg   string
	ReqID string
	Data  any
}

// 实现 error 接口
func (r *Response) Error() string { return r.Msg }

func Success(c *gin.Context, obj any, msg ...string) {
	res := &Response{
		Code: http.StatusOK,
		Data: obj,
	}
	if len(msg) > 0 {
		res.Msg = msg[0]
	}
	c.JSON(http.StatusOK, res)
}

func Fail(c *gin.Context, err error) {
	if res, ok := err.(*Response); ok {
		c.JSON(http.StatusOK, res)
		return
	}
	c.JSON(http.StatusOK, &Response{
		Code: http.StatusInternalServerError,
		Msg:  e.Cause(err).Error(),
	})
}
