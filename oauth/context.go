package oauth

import "context"

type contextKey int

const contextToken contextKey = iota

// TokenFromContext extracts token from context put there by Authenticator
func TokenFromContext(ctx context.Context) (Token, bool) {
	token, ok := ctx.Value(contextToken).(Token)
	return token, ok
}
