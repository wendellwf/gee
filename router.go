package gee

import (
	"net/http"
	"strings"
)

// roots key eg, roots['GET'], roots['POST']
// handlers key eg, handlers['GET-/p/:name/doc], handlers['POST-/p/file']

type Router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	ps := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, part := range ps {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Router) registerRouter(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *Router) searchRouter(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n == nil {
		return nil, nil
	}

	parts := parsePattern(n.pattern)
	param := make(map[string]string)
	for index, part := range parts {
		if part[0] == ':' {
			param[part[1:]] = searchParts[index]
		}
		if part[0] == '*' && len(part) > 1 {
			param[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}
	}
	return n, param
}

func (r *Router) handle(c *Context) {
	n, params := r.searchRouter(c.method, c.path)
	// fmt.Println(params)
	if n != nil {
		key := c.method + "-" + n.pattern
		c.param = params
		r.handlers[key](c)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.path)
		})
	}
	c.Next()
}
