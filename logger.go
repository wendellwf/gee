package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// start timer
		t := time.Now()
		// process request
		c.Next()
		// calculate resolution time
		log.Printf("[%d] %s in %v", c.statusCode, c.r.RequestURI, time.Since(t))
	}
}
