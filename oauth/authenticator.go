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

// NewAuthenticatorWithClient builds Authenticator with given OAuth Client
func NewAuthenticatorWithClient(c *Client) *Authenticator {
	return &Authenticator{cli: c}
}

// AuthenticateHTTP request
func (a *Authenticator) AuthenticateHTTP(ctx context.Context, kind, cred string) (context.Context, error) {
	scopes, err := a.cli.Scopes(ctx, cred)
	if err == ErrInvalidToken {
		return ctx, nil
	}

	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, contextToken, Token{Scopes: scopes}), nil
}
