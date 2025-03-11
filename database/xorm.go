package database

import (
	"time"

	"github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type DBConfig struct {
	Driver   string `yaml:"Driver"`   // 数据库驱动
	Addr     string `yaml:"Addr"`     // 数据库地址
	Port     int    `yaml:"Port"`     // 数据库端口
	UserName string `yaml:"UserName"` // 数据库用户名
	Password string `yaml:"Password"` // 数据库密码
	Database string `yaml:"Database"` // 数据库名称
}

func (cf *DBConfig) FormatDSN() string {
	config := mysql.NewConfig()
	config.User = cf.UserName
	config.Passwd = cf.Password
	config.Net = "tcp"
	config.Addr = cf.Addr
	config.DBName = cf.Database
	config.Timeout = time.Second * 5       // 连接超时
	config.ReadTimeout = time.Second * 60  // 查询超时
	config.WriteTimeout = time.Second * 60 // 插入超时
	return config.FormatDSN()
}

func NewEngine(config *DBConfig) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(config.Driver, config.FormatDSN())
	if err != nil {
		return nil, err
	}
	if err := engine.Ping(); err != nil {
		return nil, err
	}
	return engine, nil
}
