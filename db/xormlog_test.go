package db

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	xormlog "xorm.io/xorm/log"
)

func Test_xormLogger_AfterSQL(t *testing.T) {
	type args struct {
		ctx xormlog.LogContext
	}
	tests := []struct {
		name string
		l    xormlog.ContextLogger
		args args
	}{
		{
			name: "test",
			l:    newXormLogger(slog.New(slog.NewTextHandler(os.Stdout, nil))),
			args: args{
				ctx: xormlog.LogContext{
					Ctx:         context.WithValue(context.Background(), xormlog.SessionIDKey, "SessionId"),
					SQL:         "select * from user",
					ExecuteTime: time.Second * 8,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.AfterSQL(tt.args.ctx)
		})
	}
}

func Test_xormLogger_Infof(t *testing.T) {
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name string
		l    xormlog.ContextLogger
		args args
	}{
		{
			name: "test",
			l:    newXormLogger(slog.New(slog.NewTextHandler(os.Stdout, nil))),
			args: args{
				format: "test debug %s, %d",
				v:      []any{"test", 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Infof(tt.args.format, tt.args.v...)
		})
	}
}
