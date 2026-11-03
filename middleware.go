package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

// RequestIDHeader is the wire header carrying the correlation ID minted by the gateway.
const RequestIDHeader = "X-Request-ID"

// requestIDKey is the gin context key the middleware stores the ID under.
const requestIDKey = "request_id"

// RequestIDMiddleware echoes the incoming X-Request-ID and logs one line per request.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		c.Set(requestIDKey, id)
		if id != "" {
			c.Header(RequestIDHeader, id)
		}
		c.Next()
		log.Printf("request_id=%s method=%s path=%s status=%d", id, c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}

// RequestID returns the correlation ID of the current request, or "" when the middleware did not run.
func RequestID(c *gin.Context) string {
	return c.GetString(requestIDKey)
}
