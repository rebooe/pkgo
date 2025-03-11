package database

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

const sessionKey = "db/sessionKey"

func WithDB(engine xorm.EngineInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := engine.NewSession().Context(c)
		defer session.Close()

		// 开启事务
		if err := session.Begin(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Set(sessionKey, session)

		c.Next()

		// 提交事务
		if err := session.Commit(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}

func GetDB(c context.Context) xorm.Interface {
	return c.Value(sessionKey).(xorm.Interface)
}
