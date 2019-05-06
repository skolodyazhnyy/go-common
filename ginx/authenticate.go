package ginx

import (
	"context"
	"github.com/gin-gonic/gin"
	"strings"
)

// authenticator is a simple authentication service which takes context, type (kind) and credentials and returns
// new context with authentication token or error
type authenticator interface {
	AuthenticateHTTP(ctx context.Context, kind, cred string) (context.Context, error)
}

// Authenticate middleware allows to add authentication information into request context
func Authenticate(a authenticator, log logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if match := strings.Split(c.Request.Header.Get("Authorization"), " "); len(match) > 1 {
			ctx, err := a.AuthenticateHTTP(c.Request.Context(), match[0], match[1])
			if err != nil {
				log.Error("Unable to authenticate request, authenticator returned an error", map[string]interface{}{
					"context": c.Request.Context(),
					"error":   err.Error(),
				})
			} else {
				c.Request = c.Request.WithContext(ctx)
			}
		}

		c.Next()
	}
}
