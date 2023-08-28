package gins

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FileReader 返回文件
func FileReader(c *gin.Context, filename string, reader io.Reader) {
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	contentLength := 0
	if buf, ok := reader.(*bytes.Buffer); ok {
		contentLength = buf.Len()
	}
	c.DataFromReader(http.StatusOK, int64(contentLength), "application/octet-stream", reader, nil)
}
