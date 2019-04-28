package oauth

import "context"

// NewClientCredentials constructs OAuth client credential provider
// For convenience in case endpoint is not configured this provider always returns empty set of credentials
func NewClientCredentials(endpoint, client, secret string, opts ...ClientOption) *ClientCredentials {
	return &ClientCredentials{
		cli:    NewClient(endpoint, opts...),
		id:     client,
		secret: secret,
	}
}

// ClientCredentials provider
type ClientCredentials struct {
	cli    *Client
	id     string
	secret string
}

// Credentials returns Authorization header
func (c *ClientCredentials) Credentials(ctx context.Context) (string, error) {
	if c.cli == nil {
		return "", nil
	}

	token, err := c.cli.ClientCredentials(ctx, c.id, c.secret)
	if err != nil {
		return "", err
	}

	return "Bearer " + token, nil
}
