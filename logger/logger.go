package logger

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	logKey   = "logger/logkey"
	reqIdKey = "logger/reqIdKey"
)

// 获取请求id
func GetReqId(c *gin.Context) string {
	reqId := c.GetString(reqIdKey)
	if reqId == "" {
		reqId = strconv.FormatInt(time.Now().UnixNano(), 36)
		c.Set(reqIdKey, reqId)
	}
	return reqId
}

func SetLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logWith := log.With("RID", GetReqId(c))
		c.Set(logKey, logWith)
	}
}

func GetLogger(c context.Context) *slog.Logger {
	return c.Value(logKey).(*slog.Logger)
}

// 记录请求信息和错误信息
// 输出的 logger 在上下文中动态获取
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 在 defer 中调用保证在 panic 时也可以输出
		defer func() {
			logger := GetLogger(c)
			// unescape 解码请求参数
			form, _ := url.QueryUnescape(c.Request.Form.Encode())

			// 记录请求信息
			logger.Info("Request: ",
				slog.String("IP", c.ClientIP()),
				slog.String("Method", c.Request.Method),
				slog.String("Path", c.Request.URL.Path),
				slog.String("Form", form),
			)

			// 记录上下文中的错误
			if len(c.Errors) != 0 {
				for _, e := range c.Errors {
					logger.Error(fmt.Sprintf("Error: %s", e))
				}
			}
		}()

		c.Next()
	}
}
