package middle

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	logKey   = "logger/logkey"
	reqIdKey = "logger/reqIdKey"
)

func SetReqId() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := strconv.FormatInt(time.Now().UnixNano(), 36)
		c.Set(reqIdKey, reqId)
	}
}

func GetReqId(ctx context.Context) string {
	return ctx.Value(reqIdKey).(string)
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
func Logger(loger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 在 defer 中调用保证在 panic 时也可以输出
		defer func() {
			if rid := GetReqId(c); rid != "" {
				loger = loger.With("RID", rid)
			}

			// unescape 解码请求参数
			form, _ := url.QueryUnescape(c.Request.Form.Encode())
			// 记录请求信息
			loger.Info("Request",
				slog.String("IP", c.ClientIP()),
				slog.String("Method", c.Request.Method),
				slog.String("Path", c.Request.URL.Path),
				slog.String("Form", form),
			)

			// 记录上下文中的错误
			if len(c.Errors) != 0 {
				for _, e := range c.Errors {
					loger.Error(e.Error())
				}
			}
		}()

		c.Next()
	}
}
