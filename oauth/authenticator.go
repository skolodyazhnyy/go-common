package oauth

import (
	"context"
	"time"
)

// Authenticator authenticates HTTP request and adds token to context
type Authenticator struct {
	cli   *Client
	cache cache
}

// NewAuthenticator builds Authenticator and OAuth Client
func NewAuthenticator(url string, opts ...Option) *Authenticator {
	return NewAuthenticatorWithClient(NewClient(url, opts...), opts...)
}

// NewAuthenticatorWithClient builds Authenticator with given OAuth Client
func NewAuthenticatorWithClient(c *Client, opts ...Option) *Authenticator {
	o := &options{
		cache: &simpleCache{
			values:        make(map[string]interface{}),
			expire:        make(map[string]time.Time),
			lastSweep:     time.Now(),
			sweepInterval: time.Minute,
		},
	}

	for _, opt := range opts {
		opt(o)
	}

	return &Authenticator{cli: c, cache: o.cache}
}

// AuthenticateHTTP request
func (a *Authenticator) AuthenticateHTTP(ctx context.Context, kind, cred string) (context.Context, error) {
	scopes, err := a.scopes(ctx, cred)
	if err == ErrInvalidToken {
		return ctx, nil
	}

	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, contextToken, Token{Scopes: scopes}), nil
}

func (a *Authenticator) scopes(ctx context.Context, cred string) (scopes []string, err error) {
	key := "oauth_token_scopes_" + cred

	if a.cache.ShouldGet(key, &scopes) {
		return
	}

	scopes, err = a.cli.Scopes(ctx, cred)
	if err != nil {
		return
	}

	a.cache.ShouldSet(key, scopes, time.Minute)
	return
}
