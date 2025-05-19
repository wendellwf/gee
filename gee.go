package gee

import "net/http"

type HandlerFunc func(c *Context)

type Engine struct {
	router *Router
}

func New() *Engine {
	return &Engine{
		router: NewRouter(),
	}
}

func (e *Engine) Register(method string, pattern string, handler HandlerFunc) {
	e.router.registerRouter(method, pattern, handler)
}

func (e *Engine) Get(pattern string, handler HandlerFunc) {
	e.Register("GET", pattern, handler)
}

func (e *Engine) Post(pattern string, handler HandlerFunc) {
	e.Register("POST", pattern, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	e.router.handle(c)
}
