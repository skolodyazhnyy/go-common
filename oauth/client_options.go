package oauth

import "time"

type options struct {
	cache  cache
	client httpClient
}

// Option allows to specify additional OAuth client parameters
// See With* functions below for list of available options
type Option func(*options)

// WithClient allows to set HTTP client instance to be used to make calls to OAuth server
func WithClient(client httpClient) Option {
	return func(o *options) {
		o.client = client
	}
}

// WithCache allows to specify cache to be used to cache OAuth Scopes
func WithCache(ch cache) Option {
	return func(o *options) {
		o.cache = ch
	}
}

// WithoutCache allows to disable cache for OAuth scopes
func WithoutCache() Option {
	return func(o *options) {
		o.cache = &nopCache{}
	}
}

// WithSimpleCache will use simple in-memory cache for OAuth scopes
// Cache will be cleared every `sweep` interval
func WithSimpleCache(sweep time.Duration) Option {
	return func(o *options) {
		o.cache = &simpleCache{
			values:        make(map[string]interface{}),
			expire:        make(map[string]time.Time),
			lastSweep:     time.Now(),
			sweepInterval: sweep,
		}
	}
}
