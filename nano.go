package nanogo

import (
	"log"
	"net/http"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request)
type router struct {
	handleFuncMap map[string]HandleFunc
}

func (r *router) Add(name string, handleFunc HandleFunc) {
	r.handleFuncMap[name] = handleFunc
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{
			handleFuncMap: make(map[string]HandleFunc),
		},
	}
}
func (e *Engine) Run() {
	// 加载路由
	for key, value := range e.handleFuncMap {
		http.HandleFunc(key, value)
	}
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
