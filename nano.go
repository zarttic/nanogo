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
	name             string                            // 路由器组名称
	handleFuncMap    map[string]map[string]HandlerFunc // 处理函数映射表，键为路由路径，值为处理函数
	handlerMethodMap map[string][]string               // 处理方法映射表，键为路由路径，值为处理方法
	treeNode         *treeNode                         // 路由树节点
	middleWares      []MiddlewareFunc                  // 前置中间件函数列表
}

func (r *routerGroup) Use(middleware ...MiddlewareFunc) {
	r.middleWares = append(r.middleWares, middleware...)
}
func (r *routerGroup) methodHandle(handle HandlerFunc, ctx *Context) {
	if r.middleWares != nil {
		for _, middlewareFunc := range r.middleWares {
			handle = middlewareFunc(handle)
		}
	}
	handle(ctx)

}

func (r *routerGroup) handle(name, method string, handleFunc HandlerFunc) {
	_, ok := r.handleFuncMap[name]
	if !ok {
		r.handleFuncMap[name] = make(map[string]HandlerFunc)
	}
	_, ok = r.handleFuncMap[name][method]
	if ok {
		panic("有重复的路由")
	}
	r.handleFuncMap[name][method] = handleFunc
	r.treeNode.Put(name)
}

// ANY  任何方式
// ANY 添加一个任何请求方法的路由
func (r *routerGroup) ANY(name string, handleFunc HandlerFunc) {
	r.handle(name, ANY, handleFunc)
}

// GET 添加一个GET请求方法的路由
func (r *routerGroup) GET(name string, handleFunc HandlerFunc) {
	r.handle(name, http.MethodGet, handleFunc)
}

// POST 添加一个POST请求方法的路由
func (r *routerGroup) POST(name string, handleFunc HandlerFunc) {
	r.handle(name, http.MethodPost, handleFunc)
}

// DELETE 添加一个DELETE请求方法的路由
func (r *routerGroup) DELETE(name string, handleFunc HandlerFunc) {
	r.handle(name, http.MethodDelete, handleFunc)
}

// PUT  添加一个PUT请求方法的路由
func (r *routerGroup) PUT(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPut, handlerFunc)
}

// PATCH  添加一个PATCH请求方法的路由
func (r *routerGroup) PATCH(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPatch, handlerFunc)
}

// OPTIONS  添加一个OPTIONS请求方法的路由
func (r *routerGroup) OPTIONS(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodOptions, handlerFunc)
}

// Head 添加一个HEAD请求方法的路由
func (r *routerGroup) Head(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodHead, handlerFunc)
}

// 路由
type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	group := &routerGroup{
		name:             name,
		handleFuncMap:    make(map[string]map[string]HandlerFunc),
		handlerMethodMap: make(map[string][]string),
		treeNode: &treeNode{
			name:     "/",
			children: make([]*treeNode, 0),
		},
	}
	r.routerGroups = append(r.routerGroups, group)
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
		if node == nil || !node.isEnd {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintln(w, r.RequestURI+" not found")
			return
		}
		if node != nil {
			ctx := &Context{W: w, R: r}
			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				group.methodHandle(handle, ctx)
				return
			}
			//method 匹配
			handle, ok = group.handleFuncMap[node.routerName][method]
			if ok {
				group.methodHandle(handle, ctx)
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
