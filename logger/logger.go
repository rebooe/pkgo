package logger

import (
	"context"
	"log/slog"
	"project/pkg/gins"

	"github.com/gin-gonic/gin"
)

const logKey = "logger/logkey"

func SetLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logWith := log.With("RID", gins.GetReqId(c))
		c.Set(logKey, logWith)
	}
}

func GetLogger(c context.Context) *slog.Logger {
	return c.Value(logKey).(*slog.Logger)
}
