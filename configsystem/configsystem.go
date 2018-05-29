package configsystem

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	http2 "github.com/magento-mcom/go-common/configsystem/http"
	"github.com/magento-mcom/go-common/configsystem/structs"
	"time"
)

const CLIENT_LIST_URL string = "/client"

const CLIENT_SCOPES_URL string = "/client/%s/environment/%s/scope"

const CLIENT_ENVIRONMENT_URL string = "/client/%s/environment/%s/"

const SCHEMA_URL string = "/schema/"

const MERGED_URL string = "/client/%s/environment/%s/scope/%s/merged"

const LAST_MODIFIED_HEADER string = "Last-Modified"

const TIME_FORMAT string = "Mon, 02 Jan 2006 15:04:05 MST"

type ConfigSystemInterface interface {
	GetOmsClients(configurationSystemClient http2.HttpClientInterface) ([]structs.Client, error)
	GetClientScopes(configurationSystemClient http2.HttpClientInterface, client string) (structs.Scope, error)
	GetClientEnvironment(configurationSystemClient http2.HttpClientInterface, client string) (structs.Environment, error)
	HeadClientEnvironment(configurationSystemClient http2.HttpClientInterface, client string)  (structs.Environment, error)
	HeadSchema(configurationSystemClient http2.HttpClientInterface)  (structs.Schema, error)
	GetMerged(configurationSystemClient http2.HttpClientInterface, client string, scope string) (string, error)
}

type ConfigSystem struct {
	Environment string
	Url         string
}

func NewConfigSystem(environment string, url string) *ConfigSystem {
	return &ConfigSystem{
		environment,
		url,
	}
}

func (config ConfigSystem) GetOmsClients(configurationSystemClient http2.HttpClientInterface) ([]structs.Client, error) {

	url := config.Url + CLIENT_LIST_URL

	resp, err := configurationSystemClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("(url: %s, status: %d)", url, resp.StatusCode)
	}

	omsClients := make([]structs.Client, 0)
	json.NewDecoder(resp.Body).Decode(&omsClients)

	return omsClients, nil
}

func (config ConfigSystem) GetClientScopes(configurationSystemClient http2.HttpClientInterface, client string) (structs.Scope, error) {

	url := config.Url + fmt.Sprintf(CLIENT_SCOPES_URL, client, config.Environment)

	resp, err := configurationSystemClient.Get(url)

	if err != nil {
		return structs.Scope{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return structs.Scope{}, fmt.Errorf("(url: %s, status: %d)", url, resp.StatusCode)
	}

	defer resp.Body.Close()

	clientScope := structs.Scope{}
	json.NewDecoder(resp.Body).Decode(&clientScope)

	return clientScope, nil
}

func (config ConfigSystem) GetClientEnvironment(configurationSystemClient http2.HttpClientInterface, client string) (structs.Environment, error) {

	url := config.Url + fmt.Sprintf(CLIENT_ENVIRONMENT_URL, client, config.Environment)

	resp, err := configurationSystemClient.Get(url)

	if err != nil {
		return structs.Environment{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return structs.Environment{}, fmt.Errorf("(url: %s, status: %d)", url, resp.StatusCode)
	}

	defer resp.Body.Close()

	environment := structs.Environment{}
	json.NewDecoder(resp.Body).Decode(&environment)

	return environment, nil
}

func (config ConfigSystem) HeadClientEnvironment(configurationSystemClient http2.HttpClientInterface, client string)  (structs.Environment, error) {

	url := config.Url + fmt.Sprintf(CLIENT_ENVIRONMENT_URL, client, config.Environment)

	resp, err := configurationSystemClient.Head(url)
	defer resp.Body.Close()

	if err != nil {
		return structs.Environment{}, err
	}

	if resp.StatusCode == http.StatusNoContent {
		lastModified := resp.Header.Get(LAST_MODIFIED_HEADER)
		t, err := time.Parse(TIME_FORMAT, lastModified)
		if err != nil {
			return structs.Environment{}, err
		}

		environment := structs.Environment{}
		environment.ID = config.Environment
		environment.LastModification = int(t.UnixNano() / int64(time.Millisecond))

		return environment, nil
	}

	return structs.Environment{}, fmt.Errorf("Client has not environment")
}

func (config ConfigSystem) HeadSchema(configurationSystemClient http2.HttpClientInterface)  (structs.Schema, error) {

	url := config.Url + SCHEMA_URL

	resp, err := configurationSystemClient.Head(url)
	defer resp.Body.Close()

	if err != nil {
		return structs.Schema{}, err
	}

	if resp.StatusCode == http.StatusNoContent {
		lastModified := resp.Header.Get(LAST_MODIFIED_HEADER)
		t, err := time.Parse(TIME_FORMAT, lastModified)
		if err != nil {
			return structs.Schema{}, err
		}

		schema := structs.Schema{}
		schema.LastModification = int(t.UnixNano() / int64(time.Millisecond))

		return schema, nil
	}

	return structs.Schema{}, fmt.Errorf("Schema has not found")
}

func (config ConfigSystem) GetMerged(configurationSystemClient http2.HttpClientInterface, client string, scope string) (string, error) {

	url := config.Url + fmt.Sprintf(MERGED_URL, client, config.Environment, scope)

	resp, err := configurationSystemClient.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("(url: %s, status: %d)", url, resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	} else {
		bodyString := string(bodyBytes)
		return bodyString, nil
	}
}
