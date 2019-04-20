package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type logger interface {
	Error(msg string, data map[string]interface{})
	Debug(msg string, data map[string]interface{})
}

// Log middleware logs every request to the server
func Log(log logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		defer func() {
			latency := time.Since(start)
			status := c.Writer.Status()

			data := map[string]interface{}{
				"http_method":    c.Request.Method,
				"http_url":       c.Request.URL.String(),
				"http_status":    status,
				"http_useragent": c.Request.UserAgent(),
				"http_referer":   c.Request.Referer(),
				"http_latency":   latency,
				"http_protocol":  c.Request.Proto,
			}

			if status/100 == 5 {
				log.Error(fmt.Sprintf("%s request to %s failed with status code %3d", method, path, status), data)
			} else {
				log.Debug(fmt.Sprintf("%s request to %s with status code %3d", method, path, status), data)
			}
		}()

		c.Next()
	}
}
