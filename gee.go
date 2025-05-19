package gee

import "net/http"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

func (e *Engine) Register(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handler
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, r)
	} else {
		w.Write([]byte("404 Not Found"))
	}
}

func (e *Engine) Get(pattern string, handler HandlerFunc) {
	e.Register("GET", pattern, handler)
}

func (e *Engine) Post(pattern string, handler HandlerFunc) {
	e.Register("POST", pattern, handler)
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}
