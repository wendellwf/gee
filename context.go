// Package gee: this is context
package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]any

type Context struct {
	// origin objects
	w http.ResponseWriter
	r *http.Request
	// request info
	method string
	path   string
	param  map[string]string
	// response info
	statusCode int
	// middleware
	handlers []HandlerFunc
	index    int
	// engine pointer
	engine *Engine
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:      w,
		r:      r,
		method: r.Method,
		path:   r.URL.Path,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Param(key string) string {
	return c.param[key]
}

func (c *Context) FormValue(key string) string {
	return c.r.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.r.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.statusCode = code
	c.w.WriteHeader(code)
}

func (c *Context) GetStatus() int {
	return c.statusCode
}

func (c *Context) GetReq() *http.Request {
	return c.r
}

func (c *Context) SetHeader(key, value string) {
	c.w.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.w.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.w)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.w, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.w.Write(data)
}

func (c *Context) HTML(code int, name string, data any) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)

	if err := c.engine.htmlTemplates.ExecuteTemplate(c.w, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

func (c *Context) Fail(code int, context string) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.w.Write([]byte(context))
}

func (c *Context) Path() string {
	return c.path
}
