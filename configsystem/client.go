package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/magento-mcom/go-common/httpx"
)

// client communicates with config-system
type client struct {
	url        string
	env        string
	httpClient httpx.Client
}

// NewClient creates a new client instance
func NewClient(url string, env string, httpClient httpx.Client) *client {
	return &client{
		url:        url,
		env:        env,
		httpClient: httpClient,
	}
}

// Value retrieves a single key
func (c *client) Value(client string, scope string, key string, v interface{}) error {
	url := escapef("%s/client/%s/environment/%s/scope/%s/merged/%s", c.url, client, c.env, scope, key)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, v)

	if err != nil {
		return err
	}

	return nil
}

func escapef(format string, a ...interface{}) string {
	s := fmt.Sprintf(format, a...)

	return url.PathEscape(s)
}
