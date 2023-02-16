## 简单实现一个 web 框架

首先看一个 go net/http 实现的例子

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
```

http.ListenAndServe("listenAddr",nil)
这里的这个 nil，指的是默认的多路复用器
http.DefaultServeMux,它可以接收请求，并对请求进行匹配，本质上是维护了一个 map

http.HandleFunc()接收两个参数，一个是 string，代表路由，一个是 http.Handler，代表了路由对应的处理方法
这个 handler 是一个接口

```go
package http

type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}

func ListenAndServe(address string, h Handler) error
```

凡是实现了 ServeHTTP 方法的实例，都是一个 handler

但是这里有一个疑问，我们自己写的方法
`func helloHandler(w http.ResponseWriter,r *http.Request)`
并没有实现 ServeHTTP(wr)方法，为什么它也是 handler 呢？

原因在于，http.HandleFunc()里面调用了 DefaultServeMux.HandleFunc()方法

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```

这个 HandleFunc()方法内部实现:

```go
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}
```

它在内部又调用了一个 mux.Handle()

这里面的 HandlerFunc()是一个函数类型
它实现了 serveHTTP()方法，它可以将具有适当签名的函数适配成一个 handler
这里其实是做了一个转换

### 1、实现我们自己的多路复用器

```go
// 这个HandlerFunc用来定义路由映射的处理方法
type HandlerFunc func(http.ResponseWriter, *http.Request)

// 引擎 维护一个router的map
// 这个map的key由请求方法和静态路由地址构成
// 例如GET-/ POST-/ DELETE-/
// 针对相同的路由 请求方法不同也可以映射到不同的处理方法
type Engine struct {
	router map[string]HandlerFunc
}

// 用户调用New得到一个引擎
func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

// 路由注册的方法
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// 当用户调用不同的(*Engine)方法，可以将路由和处理方法的映射注册到Engine的map里
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 这里简单包装一下http.ListenAndServe()方法
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 这里实现ServeHTTP方法，将注册的路由映射取出，执行处理的函数
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, r)
        return
	}
	fmt.Fprintf(w, "404 NOT FOUND:%s\n", r.URL)
}


```
