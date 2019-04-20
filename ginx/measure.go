package ginx

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type telemetry interface {
	Timing(context.Context, string, time.Duration, []string)
}

// Measure middleware reports http_request metric
func Measure(meter telemetry) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		defer func() {
			method := c.Request.Method
			statusCode := fmt.Sprint(c.Writer.Status())
			statusClass := fmt.Sprintf("%dxx", c.Writer.Status()/100)

			meter.Timing(c.Request.Context(), "http_request", time.Since(start), []string{
				"method:" + method,
				"status:" + statusCode,
				"status_class:" + statusClass,
			})
		}()

		c.Next()
	}
}
