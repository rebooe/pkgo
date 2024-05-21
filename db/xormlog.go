package db

import (
	"fmt"
	"log/slog"
	"strings"

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
	v := ctx.Ctx.Value(xormlog.SessionIDKey)
	if key, ok := v.(string); ok {
		logWith := l.log.With("RID", key)
		logWith.Info(SQL, "exectime", ctx.ExecuteTime)
		return
	}
	l.log.Info(SQL, "exectime", ctx.ExecuteTime)
}

func (l *xormLogger) Debugf(format string, v ...any) {
	l.log.Debug(fmt.Sprintf(format, v...))
}

func (l *xormLogger) Infof(format string, v ...any) {
	l.log.Info(fmt.Sprintf(format, v...))
}

func (l *xormLogger) Warnf(format string, v ...any) {
	l.log.Warn(fmt.Sprintf(format, v...))
}

func (l *xormLogger) Errorf(format string, v ...any) {
	l.log.Error(fmt.Sprintf(format, v...))
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
