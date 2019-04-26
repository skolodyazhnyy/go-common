package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrReadingResponse    = fmt.Errorf("unable to read OAuth response")
	ErrParsingResponse    = fmt.Errorf("unable to parse OAuth response")
)

// Client provides mechanism to fetch OAuth token and scopes
type Client struct {
	url    string
	client httpClient
	cache  cache
	logger logger
}

// New OAuth Client
// Argument endpoint should point to OAuth Server root URL, use With* functions to pass additional parameters to the client
func NewClient(endpoint string, opts ...ClientOption) *Client {
	cli := &Client{
		url: strings.TrimSuffix(endpoint, "/"),
		cache: &simpleCache{
			values:        make(map[string]interface{}),
			expire:        make(map[string]time.Time),
			lastSweep:     time.Now(),
			sweepInterval: time.Minute,
		},
		logger: &nopLogger{},
		client: &http.Client{
			Timeout: time.Duration(5 * time.Second),
		},
	}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}

// ClientCredentials returns OAuth token using client credentials
func (c *Client) ClientCredentials(id, secret string) (string, error) {
	query := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{id},
		"client_secret": []string{secret},
	}

	req, err := http.NewRequest(http.MethodPost, c.url+"/oauth/token", strings.NewReader(query.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return "", ErrInvalidCredentials
	}

	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("OAuth response status code %v", resp.StatusCode)
	}

	payload := struct {
		AccessToken string `json:"access_token"`
	}{}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", ErrReadingResponse
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return "", ErrParsingResponse
	}

	return payload.AccessToken, nil
}

// Scopes for OAuth token
func (c *Client) Scopes(token string) (scopes []string, err error) {
	key := "oauth_token_scopes_" + token

	ok, err := c.cache.Get(key, &scopes)
	if err == nil && ok {
		return scopes, nil
	}

	if err != nil {
		c.logger.Error("Unable to read OAuth scopes cache", map[string]interface{}{"error": err.Error()})
	}

	req, err := http.NewRequest(http.MethodGet, c.url+"/token/"+url.QueryEscape(string(token))+"/scopes", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer resp.Body.Close()

	// token is invalid
	if resp.StatusCode == 401 {
		return nil, ErrTokenInvalid
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("OAuth server returned non-200 response: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrReadingResponse
	}

	data := struct {
		Scopes []string `json:"scopes"`
	}{}

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, ErrParsingResponse
	}

	if err := c.cache.Set(key, data.Scopes, time.Minute); err != nil {
		c.logger.Error("Unable to write OAuth scopes cache", map[string]interface{}{"error": err.Error()})
	}

	return data.Scopes, nil
}
