package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *Router {
	r := NewRouter()
	r.registerRouter("GET", "/", nil)
	r.registerRouter("GET", "/hello/:name", nil)
	r.registerRouter("GET", "hello/b/c", nil)
	r.registerRouter("GET", "/hi/:name", nil)
	r.registerRouter("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRouter(t *testing.T) {
	r := newTestRouter()
	n, ps := r.searchRouter("GET", "/hello/geektutu")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("match path: %s, params['name']: %s\n", n.pattern, ps["name"])
}
