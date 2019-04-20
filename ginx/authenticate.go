package ginx

import (
	"context"
	"github.com/gin-gonic/gin"
	"strings"
)

// authenticator is a simple authentication service which takes context, username and password and returns
// new context with authentication token and boolean flag to indicate if authentication was successful
type authenticator interface {
	AuthenticateHTTP(ctx context.Context, user, password string) (context.Context, bool)
}

// Authenticate middleware allows to add authentication information into request context
func Authenticate(a authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if match := strings.Split(c.Request.Header.Get("Authorization"), " "); len(match) > 1 {
			if ctx, ok := a.AuthenticateHTTP(c.Request.Context(), match[0], match[1]); ok {
				c.Request = c.Request.WithContext(ctx)
			}
		}

		c.Next()
	}
}
