package configsystem

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Client communicates with config-system
type Client struct {
	url        string
	env        string
	httpClient httpClient
}

// NewClient creates a new client instance
func NewClient(url string, env string, httpClient httpClient) *Client {
	return &Client{
		url:        url,
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
	data := []struct {
		ID string `json:"id"`
	}{}

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

	body, err := c.ValueAsString(client, scope, key)

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(body), v)
}

// Value String retrieves a single key when the value is a string (not an object)
func (c *Client) ValueAsString(client string, scope string, key string) (string, error) {
	endpoint := escapef("/client/%s/environment/%s/scope/%s/merged/%s", client, c.env, scope, key)

	body, err := c.get(endpoint)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) get(endpoint string) ([]byte, error) {
	addr := c.url + endpoint

	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status %d requesting %s", res.StatusCode, addr)
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func escapef(format string, a ...string) string {
	var parts []interface{}

	for _, arg := range a {
		parts = append(parts, url.PathEscape(arg))
	}

	return fmt.Sprintf(format, parts...)
}
