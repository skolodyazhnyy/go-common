package httpx

import (
	"fmt"
	"net/http"
)

type recovery interface {
	Error(msg string, data map[string]interface{})
}

// Recovery middleware allows to gracefully handle panics raised during request processing
func Recovery(log recovery) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(fmt.Sprint(err), nil)
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}()

			h.ServeHTTP(rw, req)
		})
	}
}
