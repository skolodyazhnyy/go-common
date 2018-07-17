package httpx

import (
	"net"
	"net/http"
	"strings"
)

// Local middleware makes handler available only from local address, otherwise it returns 403 error.
func Local() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ip, _, _ := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr))

			switch ip {
			case "127.0.0.1", "::1":
				h.ServeHTTP(rw, req)
			default:
				rw.WriteHeader(http.StatusForbidden)
				rw.Write([]byte("forbidden"))
			}
		})
	}
}
