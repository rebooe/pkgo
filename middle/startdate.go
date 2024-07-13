package middle

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Version() gin.HandlerFunc {
	t := time.Now()
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"time": t.Format("2006-01-02 15:04:05"),
		})
	}
}
