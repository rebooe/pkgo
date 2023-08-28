package document

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

//go:embed index.html
var fs embed.FS

type Document struct {
	once sync.Once
	dirs directorys
	info map[string]documentInfo
}

// 接口文档导航栏信息
type directory struct {
	Name string // 接口名称
	Api  string // 接口地址
}

type directorys []directory

// 接口文档详细信息
type documentInfo struct {
	// 注释信息
	Name string // 接口名称
	Api  string // 接口地址
	Auth bool   // 需要认证
	Desc string // 接口说明
	// gin 记录信息
	Method string // 请求类型
	// 反射获取信息
	Request  any // 请求数据
	Response any // 响应数据
}

func NewDocument(workdir string, paths ...string) (*Document, error) {
	doc := &Document{
		dirs: make(directorys, 0),
		info: make(map[string]documentInfo),
	}

	for i := range paths {
		if err := doc.parserDir(workdir + paths[i]); err != nil {
			return nil, fmt.Errorf("%s, path = %s", err, paths[i])
		}
	}
	// 对接口排序
	sort.Slice(doc.dirs, func(i, j int) bool {
		return doc.dirs[i].Api < doc.dirs[j].Api
	})
	return doc, nil
}

// 解析目录下的文件注释信息
func (d *Document) parserDir(path string) error {
	// 匹配注释段的正则
	var commentCmp = regexp.MustCompile(`/\*:doc([\s\S]*?)\*/`)

	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// 遍历目录
	for _, dir := range dirEntry {
		// 读取文件内容
		if dir.IsDir() {
			continue
		}
		file, err := os.Open(filepath.Join(path, dir.Name()))
		if err != nil {
			return err
		}
		defer file.Close()
		fByte, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		// 正则匹配文件内的注释
		matchAll := commentCmp.FindAllSubmatch(fByte, -1)

		for _, match := range matchAll {
			// 解析每个匹配的注释段,按照*号分行
			lineScile := bytes.Split(match[1], []byte("*"))

			// 解析注释段的每一行内容到结构体
			var info documentInfo
			for _, line := range lineScile {
				line := bytes.TrimSpace(line)

				before, after, found := strings.Cut(string(line), " ")
				if !found {
					continue
				}
				switch before {
				case "@Name":
					info.Name = after
				case "@Api":
					info.Api = after
				case "@Auth":
					info.Auth = after == "true"
				case "@Description":
					info.Desc = after
				default:
					continue
				}
			}
			d.info[info.Api] = info
			// 添加导航栏信息
			d.dirs = append(d.dirs, directory{Name: info.Name, Api: info.Api})
		}
	}
	return nil
}

func (d *Document) Handler(engine *gin.Engine, handlers ...gin.HandlerFunc) {
	group := engine.Group("/api/doc", handlers...)
	group.GET("/", d.index(engine))
	group.GET("/nav", d.navigator)
	group.GET("/info", d.docinfo)
}

func (d *Document) index(engine *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取路由信息，请求时初始化可以保证在所有路由注册后执行
		d.once.Do(func() {
			for _, route := range engine.Routes() {

				if info, ok := d.info[route.Path]; ok {
					info.Method = route.Method

					// 解析 Handler 函数获取出入参数
					hanleType := reflect.TypeOf(route.HandlerFunc)
					if hanleType != nil {
						if hanleType.NumIn() > 1 {
							info.Request = hanleType.In(1)
						}

						if hanleType.NumOut() > 0 {
							info.Response = hanleType.Out(0)
						}
					}

					// 更新文档信息
					d.info[route.Path] = info
				}
			}
		})

		c.FileFromFS("index.html", http.FS(fs))
	}
}

func (d *Document) navigator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Data": d.dirs})
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

func (d *Document) docinfo(c *gin.Context) {
	q := c.Query("q")
	info, ok := d.info[q]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Msg": "查询接口不存在"})
		return
	}

	if info.Request != nil {
		rt := info.Request.(reflect.Type)
		if rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
		}
		if rt.Kind() == reflect.Struct { // 只支持结构体解析
			args := make([]infoRequest, rt.NumField())

			for i := 0; i < rt.NumField(); i++ {
				field := rt.Field(i)
				arg := infoRequest{
					Name:    field.Name,
					Type:    field.Type.Kind().String(),
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

				args[i] = arg
			}

			info.Request = args
		}
	}

	if info.Response != nil {
		rt := info.Response.(reflect.Type)
		info.Response = parserResponse(rt)
	}

	c.JSON(http.StatusOK, gin.H{"Data": info})
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
			Type:    field.Type.Kind().String(),
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
