package gins

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
)

type testForm struct {
	A string
	B int
}

func test(c *gin.Context, form testForm) *testForm {
	return nil
}

func TestWarp(t *testing.T) {
	engine := New()

	group := engine.Group("/")
	group.GET("/test", test)
	group2 := group.Group("/group2")
	group2.POST("/test", test)
	t.Logf("%v", engine.Routes())

	server := httptest.NewServer(engine)
	defer server.Close()

	var body []byte
	err := gout.GET(server.URL + "/test").
		SetQuery(gout.H{
			"A": "我是A",
			"B": 123,
		}).
		BindBody(&body).Do()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", body)
}
