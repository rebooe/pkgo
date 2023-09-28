package db

import (
	"context"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

const sessionKey = "db/sessionKey"

func SetDB(engine *xorm.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := engine.NewSession().Context(c)
		defer session.Close()

		c.Set(sessionKey, session)
		c.Next()
	}
}

func GetDB(c context.Context) *xorm.Session {
	return c.Value(sessionKey).(*xorm.Session)
}
