package oauth

import "context"

// Authenticator authenticates HTTP request and adds token to context
type Authenticator struct {
	cli *Client
}

// NewAuthenticator builds Authenticator and OAuth Client
func NewAuthenticator(url string, opts ...ClientOption) *Authenticator {
	return NewAuthenticatorWithClient(NewClient(url, opts...))
}

// NewAuthenticator builds Authenticator with given OAuth Client
func NewAuthenticatorWithClient(c *Client) *Authenticator {
	return &Authenticator{cli: c}
}

// AuthenticateHTTP request
func (a *Authenticator) AuthenticateHTTP(ctx context.Context, user, password string) (context.Context, bool) {
	scopes, err := a.cli.Scopes(ctx, password)
	if err == ErrInvalidToken {
		return ctx, false
	}

	if err != nil {
		return ctx, false
	}

	return context.WithValue(ctx, contextToken, Token{Scopes: scopes}), true
}
