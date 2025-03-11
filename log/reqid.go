package log

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	reqIdKey = "middle/reqIdKey"
	logKey   = "middle/logkey"
)

func WithReqId() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := strconv.FormatInt(time.Now().UnixNano(), 36)
		c.Set(reqIdKey, reqId)
	}
}

func GetReqId(ctx context.Context) string {
	return ctx.Value(reqIdKey).(string)
}
