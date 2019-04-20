package httpx

import (
	"fmt"
	"net/http"
	"time"
)

type telemetry interface {
	Timing(string, time.Duration, []string)
}

// Measure middleware reports http_request metric
func Measure(meter telemetry) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			writer := &responseWriter{ResponseWriter: rw, status: http.StatusOK}
			start := time.Now()

			defer func() {
				method := req.Method
				statusCode := fmt.Sprint(writer.status)
				statusClass := fmt.Sprintf("%dxx", writer.status/100)

				meter.Timing("http_request", time.Since(start), []string{
					"method:" + method,
					"status:" + statusCode,
					"status_class:" + statusClass,
				})
			}()

			h.ServeHTTP(writer, req)
		})
	}
}
