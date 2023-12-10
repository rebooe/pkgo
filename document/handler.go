package document

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

//go:embed index.html
var indexHtml embed.FS

type response struct {
	Code int
	Msg  string
	Data any
}

func (doc *Document) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		doc.index(w, r)
	case "/nav":
		doc.navigator(w, r)
	case "/info":
		doc.docinfo(w, r)
	default:
		fmt.Fprintf(w, "请求错误 %s", r.URL.Path)
	}
}

func (doc *Document) index(w http.ResponseWriter, r *http.Request) {
	file, _ := indexHtml.ReadFile("index.html")
	w.Write(file)
}

// 接口文档导航栏信息
type navigatorRes struct {
	Name string // 接口名称
	Api  string // 接口地址
}

func (doc *Document) navigator(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/json")

	res := make([]navigatorRes, 0, len(doc.info))
	for k, v := range doc.info {
		res = append(res, navigatorRes{Name: v.Name, Api: k})
	}

	encode := json.NewEncoder(w)
	encode.Encode(&response{Code: 0, Data: res})
}

type docInfoRes struct {
	Info
	Request  []infoRequest
	Response []infoResponse
}

type infoRequest struct {
	Name     string // 名称
	Type     string // 类型
	Default  string // 默认值
	Desc     string // 说明
	Required bool   // 必填
}

type infoResponse struct {
	Name    string         // 名称
	Type    string         // 类型
	Default string         // 示例
	Desc    string         // 说明
	Child   []infoResponse // 下级结构
}

func (doc *Document) docinfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/json")
	encode := json.NewEncoder(w)

	q := r.URL.Query().Get("q")
	info, ok := doc.info[q]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		encode.Encode(&response{Code: 1, Msg: "查询接口不存在"})
		return
	}

	res := docInfoRes{Info: info}

	if info.Request != nil {
		rt := info.Request
		if rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
		}
		if rt.Kind() == reflect.Struct { // 只支持结构体解析
			res.Request = make([]infoRequest, rt.NumField())

			for i := 0; i < rt.NumField(); i++ {
				field := rt.Field(i)
				arg := infoRequest{
					Name:    field.Name,
					Type:    typeString(field.Type.Kind()),
					Default: field.Tag.Get("default"),
					Desc:    field.Tag.Get("desc"),
				}
				if name, ok := field.Tag.Lookup("form"); ok {
					arg.Name = name
				}
				if field.Type.Kind() == reflect.Pointer {
					arg.Type = field.Type.Elem().Kind().String()
				}

				if strings.Contains(field.Tag.Get("binding"), "required") {
					arg.Required = true
				}

				res.Request[i] = arg
			}
		}
	}

	if info.Response != nil {
		rt := info.Response
		res.Response = parserResponse(rt)
	}

	encode.Encode(&response{Code: 0, Data: res})
}

func parserResponse(obj reflect.Type) []infoResponse {
	rt := obj
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return nil
	}

	infos := make([]infoResponse, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		info := infoResponse{
			Name:    field.Name,
			Type:    typeString(field.Type.Kind()),
			Desc:    field.Tag.Get("desc"),
			Default: field.Tag.Get("default"),
		}
		if name, ok := field.Tag.Lookup("json"); ok {
			info.Name = name
		}
		if field.Type.Kind() == reflect.Pointer {
			info.Type = field.Type.Elem().Kind().String()
		}

		info.Child = parserResponse(field.Type)
		infos[i] = info
	}
	return infos
}

func typeString(str fmt.Stringer) string {
	s := str.String()
	switch s {
	case "slice":
		return "array"
	case "struct":
		return "object"
	}
	return s
}
