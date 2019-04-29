package httpx

import (
	"context"
	"net/http"
	"strings"
)

// authenticator is a simple authentication service which takes context, kind (type) and credentials and returns
// new context with authentication token or error
type authenticator interface {
	AuthenticateHTTP(ctx context.Context, kind, cred string) (context.Context, error)
}

// Authenticate middleware allows to add authentication information into request context
func Authenticate(a authenticator, log logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if match := strings.SplitN(req.Header.Get("Authorization"), " ", 2); len(match) > 1 {
				ctx, err := a.AuthenticateHTTP(req.Context(), match[0], match[1])
				if err != nil {
					log.Error("Unable to authenticate request, authenticator returned an error", map[string]interface{}{
						"context": req.Context(),
						"error":   err.Error(),
					})
				} else {
					req = req.WithContext(ctx)
				}
			}

			h.ServeHTTP(rw, req)
		})
	}
}
