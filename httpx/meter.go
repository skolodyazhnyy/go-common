package httpx

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

// Meter middleware reports metrics to prometheus
func Meter() func(http.Handler) http.Handler {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "HTTP Request Counter",
	}, []string{"method", "path", "status_class"})

	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_requests_duration_ms",
		Help:    "HTTP Request Duration",
		Buckets: prometheus.LinearBuckets(10, 10, 10),
	}, []string{"method", "path", "status_class"})

	prometheus.MustRegister(counter, duration)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			writer := &responseWriter{ResponseWriter: rw, status: http.StatusOK}
			start := time.Now()

			defer func() {
				method := req.Method
				path := req.URL.Path
				status := fmt.Sprintf("%dxx", writer.status/100)

				duration.WithLabelValues(method, path, status).Observe(time.Since(start).Seconds())
				counter.WithLabelValues(method, path, status).Add(1)
			}()

			h.ServeHTTP(writer, req)
		})
	}
}
