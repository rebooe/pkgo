package middle

import (
	"log"
	"strings"

	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	reqIdKey  = "middle/logger/reqid_key"
	logKey    = "middle/logger/log_key"
	userIDKey = "middle/logger/userid_key"
)

// Logger 捕获程序的异常，并将错误信息输出到日志
func Logger(out *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成唯一标识符
		reqId := strconv.FormatInt(time.Now().UnixNano(), 36)
		c.Set(reqIdKey, reqId)
		req := c.Request

		defer func() {
			form, _ := url.QueryUnescape(req.Form.Encode())
			out.Printf("[INFO][%s][%s][%v] %s %s\n%s", reqId, c.ClientIP(), c.Value(userIDKey), req.Method, req.URL.Path, form)
			// 记录跟踪日志
			ss := c.GetStringSlice(logKey)
			if len(ss) > 0 {
				out.Printf("[INFO][%s]Trace:\n%s\n\n", reqId, strings.Join(ss, "\n"))
			}
		}()

		defer func() {
			if err := recover(); err != nil {
				out.Printf("[PANIC][%s]%s", reqId, err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Msg": "系统异常", "ReqId": reqId})
			}
		}()

		c.Next()

		// 记录错误日志
		if len(c.Errors) != 0 {
			msg := ""
			for _, e := range c.Errors {
				msg = fmt.Sprintf("%s%s\n", msg, e)
			}
			out.Printf("[ERROR][%s]%s", reqId, msg)
		}
	}
}

// 设置用户id
func SetUseID(c *gin.Context, id any) {
	c.Set(userIDKey, id)
}

// 设置跟踪日志输出
func SetLog(c *gin.Context, format string, args ...any) {
	ss := c.GetStringSlice(logKey)
	ss = append(ss, fmt.Sprintf(format, args...))
	c.Set(logKey, ss)
}

// 获取请求id
func GetReqId(c *gin.Context) string {
	return c.GetString(reqIdKey)
}
