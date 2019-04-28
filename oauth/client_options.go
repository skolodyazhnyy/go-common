package oauth

import "time"

// ClientOption allows to specify additional OAuth client parameters
// See With* functions below for list of available options
type ClientOption func(*Client)

// WithClient allows to set HTTP client instance to be used to make calls to OAuth server
func WithClient(client httpClient) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

// WithCache allows to specify cache to be used to cache OAuth Scopes
func WithCache(ch cache) ClientOption {
	return func(c *Client) {
		c.cache = ch
	}
}

// WithoutCache allows to disable cache for OAuth scopes
func WithoutCache() ClientOption {
	return func(c *Client) {
		c.cache = &nopCache{}
	}
}

// WithSimpleCache will use simple in-memory cache for OAuth scopes
// Cache will be cleared every `sweep` interval
func WithSimpleCache(sweep time.Duration) ClientOption {
	return func(c *Client) {
		c.cache = &simpleCache{
			values:        make(map[string]interface{}),
			expire:        make(map[string]time.Time),
			lastSweep:     time.Now(),
			sweepInterval: sweep,
		}
	}
}

// WithLogger allows to specify logger which will be used to log errors
func WithLogger(log logger) ClientOption {
	return func(c *Client) {
		c.logger = log
	}
}
