package configsystem

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPClient is well-known HTTP client interface
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Client communicates with config-system
type Client struct {
	host       string
	env        string
	httpClient HTTPClient
}

// NewClient creates a new client instance
func NewClient(host string, env string, httpClient HTTPClient) *Client {
	return &Client{
		host:       host,
		env:        env,
		httpClient: httpClient,
	}
}

// Clients retrieves a list of available clients
func (c *Client) Clients() ([]string, error) {
	endpoint := "/client"

	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	clients := []string{}
	data := []struct{ ID string }{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	for _, elem := range data {
		clients = append(clients, elem.ID)
	}

	return clients, nil
}

// Value retrieves a single key
func (c *Client) Value(client string, scope string, key string, v interface{}) error {
	endpoint := escapef("/client/%s/environment/%s/scope/%s/merged/%s", client, c.env, scope, key)

	body, err := c.get(endpoint)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.host, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status %d requesting %s", res.StatusCode, url)
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func escapef(format string, a ...string) string {
	parts := []interface{}{}

	for _, arg := range a {
		parts = append(parts, url.PathEscape(arg))
	}

	return fmt.Sprintf(format, parts...)
}
