package log

import (
	"context"
	"io"
	"log/slog"
)

type LoggerConfig struct {
	Level string `yaml:"Level"` // 日志级别
}

type Logger interface {
	Log(ctx context.Context, level slog.Level, msg string, args ...any)
}

type log struct {
	*slog.Logger
}

func NewLog(config *LoggerConfig, writer io.Writer) Logger {
	level := slog.LevelInfo
	switch config.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	return &log{
		Logger: slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level})),
	}
}

func (l *log) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	rid := GetReqId(ctx)
	if rid != "" {
		l.Logger.With(slog.String("RID", rid)).
			Log(ctx, level, msg, args...)
		return
	}
	l.Logger.Log(ctx, level, msg, args...)
}
