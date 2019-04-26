package oauth

import "context"

type contextKey int

const contextToken contextKey = iota

func TokenFromContext(ctx context.Context) (Token, bool) {
	token, ok := ctx.Value(contextToken).(Token)
	return token, ok
}
