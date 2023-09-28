package db

import (
	"log/slog"

	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

func New(dataSourceName string, log *slog.Logger) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	if log != nil {
		engine.SetLogger(newXormLogger(log)) // 设置输出日志
	}

	engine.SetMapper(names.GonicMapper{}) // 设置名称映射规则
	return engine, nil
}
