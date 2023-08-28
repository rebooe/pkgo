package gins

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type RouteInfo struct {
	Method      string
	Path        string
	HandlerFunc any
}

func WarpAny(handler any) gin.HandlerFunc {
	// 反射获取函数类型
	handlerValue := reflect.ValueOf(handler)
	handleType := handlerValue.Type()
	if handleType.Kind() != reflect.Func {
		panic("HandlerFunc 类型必须是函数")
	}
	// 第二个参数为请求表单模型(如果有)
	var formType reflect.Type
	if handleType.NumIn() > 1 {
		formType = handleType.In(1)
	}

	return func(c *gin.Context) {
		// 准备函数参数
		args := []reflect.Value{reflect.ValueOf(c)}

		// 验证输入表单
		if formType != nil {
			form := reflect.New(formType).Interface()
			if err := c.Bind(form); err != nil { // TODO: 1.会自动调用c.Error()
				return
			}
			args = append(args, reflect.ValueOf(form).Elem())
		}

		handlerValue.Call(args)
	}
}

type Engine struct {
	*gin.Engine
	routeInfos []RouteInfo
}

func New() *Engine {
	return &Engine{
		Engine:     gin.New(),
		routeInfos: []RouteInfo{},
	}
}

func (engine *Engine) Routes() []RouteInfo {
	return engine.routeInfos
}

func (engine *Engine) Group(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	return &RouterGroup{
		RouterGroup: engine.RouterGroup.Group(relativePath, handlers...),
		engine:      engine,
	}
}

type RouterGroup struct {
	*gin.RouterGroup
	engine *Engine
}

func (route *RouterGroup) Group(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	return &RouterGroup{
		RouterGroup: route.RouterGroup.Group(relativePath, handlers...),
		engine:      route.engine,
	}
}

func (group *RouterGroup) GET(relativePath string, handlers ...any) {
	group.Handle(http.MethodGet, relativePath, handlers...)
}

func (group *RouterGroup) POST(relativePath string, handlers ...any) {
	group.Handle(http.MethodPost, relativePath, handlers...)
}

func (group *RouterGroup) Handle(httpMethod string, relativePath string, handlers ...any) gin.IRoutes {
	hands := make([]gin.HandlerFunc, len(handlers))
	routeInfo := RouteInfo{
		Method: httpMethod,
		Path:   strings.ReplaceAll(group.BasePath()+relativePath, "//", "/"),
	}

	for i := range handlers {
		switch h := handlers[i].(type) {
		case gin.HandlerFunc:
			hands[i] = h
		default:
			// 包装自定义方法
			hands[i] = WarpAny(h)
			routeInfo.HandlerFunc = handlers[i]
		}
	}
	// 记录路由相关信息
	group.engine.routeInfos = append(group.engine.routeInfos, routeInfo)

	return group.RouterGroup.Handle(httpMethod, relativePath, hands...)
}
