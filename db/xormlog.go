package db

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rebooe/pkgo/logger"

	xormlog "xorm.io/xorm/log"
)

type xormLogger struct {
	log      *slog.Logger
	logLevel xormlog.LogLevel
	showSQL  bool
}

func newXormLogger(l *slog.Logger) xormlog.ContextLogger {
	return &xormLogger{
		log:      l,
		logLevel: xormlog.LOG_INFO,
		showSQL:  true,
	}
}

func (l *xormLogger) BeforeSQL(ctx xormlog.LogContext) {}

func (l *xormLogger) AfterSQL(ctx xormlog.LogContext) {
	SQL := ReplacePlaceholders(ctx.SQL, ctx.Args...)
	if gctx, ok := ctx.Ctx.(*gin.Context); ok {
		reqId := logger.GetReqId(gctx)
		l.Infof("[%s][%v]%s", reqId, ctx.ExecuteTime, SQL)
		return
	}
	l.Infof("[%v]%s", ctx.ExecuteTime, SQL)
}

func (l *xormLogger) Debugf(format string, v ...any) {
	l.log.Debug(format, v...)
}

func (l *xormLogger) Infof(format string, v ...any) {
	l.log.Info(format, v...)
}

func (l *xormLogger) Warnf(format string, v ...any) {
	l.log.Warn(format, v...)
}

func (l *xormLogger) Errorf(format string, v ...any) {
	l.log.Error(format, v...)
}

func (l *xormLogger) Level() xormlog.LogLevel         { return l.logLevel }
func (l *xormLogger) SetLevel(level xormlog.LogLevel) { l.logLevel = level }

func (l *xormLogger) ShowSQL(show ...bool) { l.showSQL = show[0] }
func (l *xormLogger) IsShowSQL() bool      { return l.showSQL }

// ReplacePlaceholders 替换查询字符串中的 ? 为实际参数值
func ReplacePlaceholders(query string, args ...interface{}) string {
	query = strings.ReplaceAll(query, "?", "%v")
	return fmt.Sprintf(query, args...)
}
