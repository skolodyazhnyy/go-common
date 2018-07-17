package httpx

import (
	"context"
	"net/http"
	"strings"
)

// authenticator is a simple authentication service which takes context, username and password and returns
// new context with authentication token and boolean flag to indicate if authentication was successful
type authenticator interface {
	AuthenticateHTTP(ctx context.Context, user, password string) (context.Context, bool)
}

// Authenticate middleware allows to add authentication information into request context
func Authenticate(a authenticator) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if match := strings.Split(req.Header.Get("Authorization"), " "); len(match) > 1 {
				if ctx, ok := a.AuthenticateHTTP(req.Context(), match[0], match[1]); ok {
					req = req.WithContext(ctx)
				}
			}

			h.ServeHTTP(rw, req)
		})
	}
}
