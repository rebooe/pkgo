package main

import (
	"log"
	"net/http"

	"github.com/rebooe/pkg-go/document"
)

func main() {
	doc := document.NewDocument(true)

	doc.Comments(
		"/test1",
		document.Name("测试1"),
		document.Method("POST"),
		document.Desc("这是一个测试接口"),
		document.Request(struct {
			A string `default:"abc" desc:"参数A"`
		}{}),
		document.Response(struct {
			B int `default:"123" desc:"上大分"`
		}{}),
	)

	doc.Comments(
		"/test2",
		document.Name("测试2"),
		document.Method("GET"),
		document.Desc(`这是第二个测试接口`),
		document.Request(struct {
			A string `default:"abc" desc:"拉拉"`
		}{}),
		document.Response(struct {
			B int `default:"123" desc:"上大分"`
		}{}),
	)

	http.Handle("/", doc)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
