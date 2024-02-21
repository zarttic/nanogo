package nanogo

import (
	"fmt"
	"log"
	"net/http"
)

const ANY = "ANY"

// HandlerFunc 处理器
type HandlerFunc func(ctx *Context)

// 路由组
type routerGroup struct {
	name             string
	handleFuncMap    map[string]map[string]HandlerFunc
	handlerMethodMap map[string][]string
	treeNode         *treeNode
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
	method := r.Method
	for _, group := range e.routerGroups {
		routerName := SubstringLast(r.RequestURI, "/"+group.name)
		node := group.treeNode.Get(routerName)
		if node == nil || !node.isEnd {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, r.RequestURI+" not found")
			return
		}
		if node != nil {
			ctx := &Context{W: w, R: r}
			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				handle(ctx)
				return
			}
			//method 匹配
			handle, ok = group.handleFuncMap[node.routerName][method]
			if ok {
				handle(ctx)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = fmt.Fprintf(w, "[%s] %s not allowed\n", method, r.RequestURI)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s  not found\n", r.RequestURI)
}
func (e *Engine) Run() {
	http.Handle("/", e)
	err := http.ListenAndServe(":9421", nil)
	if err != nil {
		log.Fatal(err)
	}
}
