package nanogo

import (
	"fmt"
	"log"
	"net/http"
)

const ANY = "ANY"

// HandlerFunc 处理器
type HandlerFunc func(ctx *Context)

// MiddlewareFunc 中间件
type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

// routerGroup 结构体表示一个路由器组
type routerGroup struct {
	name               string                            // 路由器组名称
	handleFuncMap      map[string]map[string]HandlerFunc // 处理函数映射表，键为路由路径，值为处理函数
	middlewaresFuncMap map[string]map[string][]MiddlewareFunc
	handlerMethodMap   map[string][]string // 处理方法映射表，键为路由路径，值为处理方法
	treeNode           *treeNode           // 路由树节点
	middleWares        []MiddlewareFunc    // 前置中间件函数列表
}

func (r *routerGroup) Use(middleware ...MiddlewareFunc) {
	r.middleWares = append(r.middleWares, middleware...)
}
func (r *routerGroup) methodHandle(name string, method string, handle HandlerFunc, ctx *Context) {
	//通用
	if r.middleWares != nil {
		for _, middlewareFunc := range r.middleWares {
			handle = middlewareFunc(handle)
		}
	}
	//路由级别
	middlewareFuncs := r.middlewaresFuncMap[name][method]
	if middlewareFuncs != nil {
		for _, middlewareFunc := range middlewareFuncs {
			handle = middlewareFunc(handle)
		}
	}

	handle(ctx)

}

func (r *routerGroup) handle(name, method string, handleFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	_, ok := r.handleFuncMap[name]
	if !ok {
		r.handleFuncMap[name] = make(map[string]HandlerFunc)
		r.middlewaresFuncMap[name] = make(map[string][]MiddlewareFunc)
	}
	_, ok = r.handleFuncMap[name][method]
	if ok {
		panic("有重复的路由")
	}
	r.handleFuncMap[name][method] = handleFunc
	r.middlewaresFuncMap[name][method] = append(r.middlewaresFuncMap[name][method], middlewareFunc...)
	r.treeNode.Put(name)
}

// ANY  任何方式
// ANY 添加一个任何请求方法的路由
func (r *routerGroup) ANY(name string, handleFunc HandlerFunc, middleware ...MiddlewareFunc) {
	r.handle(name, ANY, handleFunc, middleware...)
}

// GET 添加一个GET请求方法的路由
func (r *routerGroup) GET(name string, handleFunc HandlerFunc, middleware ...MiddlewareFunc) {
	r.handle(name, http.MethodGet, handleFunc, middleware...)
}

// POST 添加一个POST请求方法的路由
func (r *routerGroup) POST(name string, handleFunc HandlerFunc, middleware ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handleFunc, middleware...)
}

// DELETE 添加一个DELETE请求方法的路由
func (r *routerGroup) DELETE(name string, handleFunc HandlerFunc, middleware ...MiddlewareFunc) {
	r.handle(name, http.MethodDelete, handleFunc, middleware...)
}

// PUT  添加一个PUT请求方法的路由
func (r *routerGroup) PUT(name string, handlerFunc HandlerFunc, middleware ...MiddlewareFunc) {
	r.handle(name, http.MethodPut, handlerFunc, middleware...)
}

// PATCH  添加一个PATCH请求方法的路由
func (r *routerGroup) PATCH(name string, handlerFunc HandlerFunc, middleware ...MiddlewareFunc) {
	r.handle(name, http.MethodPatch, handlerFunc, middleware...)
}

// OPTIONS  添加一个OPTIONS请求方法的路由
func (r *routerGroup) OPTIONS(name string, handlerFunc HandlerFunc, middleware ...MiddlewareFunc) {
	r.handle(name, http.MethodOptions, handlerFunc, middleware...)
}

// Head 添加一个HEAD请求方法的路由
func (r *routerGroup) Head(name string, handlerFunc HandlerFunc, middleware ...MiddlewareFunc) {
	r.handle(name, http.MethodHead, handlerFunc, middleware...)
}

// 路由
type router struct {
	routerGroups []*routerGroup
}

// Group Group函数用于创建一个路由组
func (r *router) Group(name string) *routerGroup {
	// 创建一个routerGroup对象
	group := &routerGroup{
		name:             name,
		handleFuncMap:    make(map[string]map[string]HandlerFunc),
		handlerMethodMap: make(map[string][]string),
		treeNode: &treeNode{
			name:     "/",
			children: make([]*treeNode, 0),
		},
		middlewaresFuncMap: make(map[string]map[string][]MiddlewareFunc),
	}
	// 将routerGroup对象添加到router的routerGroups切片中
	r.routerGroups = append(r.routerGroups, group)
	// 返回创建的routerGroup对象
	return group
}

// Add 添加路由

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{},
	}
}
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.httpRequestHandle(w, r)
}
func (e *Engine) httpRequestHandle(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	for _, group := range e.routerGroups {
		routerName := SubstringLast(r.RequestURI, "/"+group.name)
		node := group.treeNode.Get(routerName)
		if node != nil && node.isEnd {
			ctx := &Context{W: w, R: r}
			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				group.methodHandle(node.routerName, ANY, handle, ctx)
				return
			}
			//method 匹配
			handle, ok = group.handleFuncMap[node.routerName][method]
			if ok {
				group.methodHandle(node.routerName, method, handle, ctx)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = fmt.Fprintf(w, "[%s] %s not allowed\n", method, r.RequestURI)
			return
		}

	}
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprintf(w, "%s  not found\n", r.RequestURI)
}

func (e *Engine) Run() {
	http.Handle("/", e)
	err := http.ListenAndServe(":9421", nil)
	if err != nil {
		log.Fatal(err)
	}
}
