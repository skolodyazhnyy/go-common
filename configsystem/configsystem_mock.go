package configsystem

import (
	"github.com/magento-mcom/go-common/configsystem/http"
	"github.com/magento-mcom/go-common/configsystem/structs"
	"github.com/stretchr/testify/mock"
)

type ConfigSystemMock struct {
	mock.Mock
}

func (config *ConfigSystemMock) GetOmsClients(configurationSystemClient http.HttpClientInterface) ([]structs.Client, error) {
	args := config.Called(configurationSystemClient)
	return args.Get(0).([]structs.Client), args.Error(1)
}

func (config *ConfigSystemMock) GetClientScopes(configurationSystemClient http.HttpClientInterface, client string) (structs.Scope, error) {
	args := config.Called(configurationSystemClient, client)
	return args.Get(0).(structs.Scope), args.Error(1)
}

func (config *ConfigSystemMock) GetClientEnvironment(configurationSystemClient http.HttpClientInterface, client string) (structs.Environment, error) {
	args := config.Called(configurationSystemClient, client)
	return args.Get(0).(structs.Environment), args.Error(1)
}

func (config *ConfigSystemMock) HeadClientEnvironment(configurationSystemClient http.HttpClientInterface, client string) (structs.Environment, error) {
	args := config.Called(configurationSystemClient, client)
	return args.Get(0).(structs.Environment), args.Error(1)
}

func (config *ConfigSystemMock) HeadSchema(configurationSystemClient http.HttpClientInterface) (structs.Schema, error) {
	args := config.Called(configurationSystemClient)
	return args.Get(0).(structs.Schema), args.Error(1)
}

func (config *ConfigSystemMock) GetMerged(configurationSystemClient http.HttpClientInterface, client string, scope string) (string, error) {
	args := config.Called(configurationSystemClient, client, scope)
	return args.String(0), args.Error(1)
}
