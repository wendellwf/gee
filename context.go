package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

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
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:      w,
		r:      r,
		method: r.Method,
		path:   r.URL.Path,
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.param[key]
	return value
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

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.w.Write([]byte(html))
}

func (c *Context) Path() string {
	return c.path
}
