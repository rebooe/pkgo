// 请求结构体中的tag：desc:注释，default:默认值，form:请求字段名
// 响应结构体中的tag：desc:注释，default:示例值
package document

import "reflect"

type Document struct {
	open bool
	info map[string]Info
}

// 接口文档详细信息
type Info struct {
	// 注释信息
	Name string // 接口名称
	Api  string // 接口地址
	Auth bool   // 需要认证
	Desc string // 接口说明
	// gin 记录信息
	Method string // 请求类型
	// 反射获取信息
	Request  reflect.Type // 请求数据
	Response reflect.Type // 响应数据
}

func NewDocument(open bool) *Document {
	doc := &Document{
		open: open,
		info: make(map[string]Info),
	}

	// 对接口排序
	// sort.Slice(doc.dirs, func(i, j int) bool {
	// 	return doc.dirs[i].Api < doc.dirs[j].Api
	// })
	return doc
}

type WithDocument func(*Info)

func (doc *Document) Comments(api string, ops ...WithDocument) {
	if !doc.open {
		return
	}
	info := Info{Api: api, Auth: true}
	for _, op := range ops {
		op(&info)
	}
	doc.info[api] = info
}

func Name(name string) WithDocument {
	return func(i *Info) {
		i.Name = name
	}
}

func Auth(isAuth bool) WithDocument {
	return func(i *Info) {
		i.Auth = isAuth
	}
}

func Method(method string) WithDocument {
	return func(i *Info) {
		i.Method = method
	}
}

func Desc(desc string) WithDocument {
	return func(i *Info) {
		i.Desc = desc
	}
}

func Request(v any) WithDocument {
	return func(i *Info) {
		i.Request = reflect.TypeOf(v)
	}
}

func Response(v any) WithDocument {
	return func(i *Info) {
		i.Response = reflect.TypeOf(v)
	}
}
