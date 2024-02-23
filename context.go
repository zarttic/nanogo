package nanogo

import "net/http"

// Context 结构体表示一个 HTTP 请求的上下文
type Context struct {
	W http.ResponseWriter // ResponseWriter 接口用于向客户端发送 HTTP 响应
	R *http.Request       // *http.Request 表示一个 HTTP 请求
}
