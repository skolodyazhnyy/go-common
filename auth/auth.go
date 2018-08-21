package auth

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Auth struct {
	Environment string
	Url         string
	Client      string
	Secret      string
}

func NewAuth(environment string, url string, client string, secret string) *Auth {
	return &Auth{
		environment,
		url,
		client,
		secret,
	}
}

func (auth *Auth) GetAuthorizedClient() (clientcredentials.Config, error) {
	if auth.Client == "" || auth.Secret == "" {
		return clientcredentials.Config{}, fmt.Errorf("Configuration parameters for Auth must be set")
	}
	oauth2.RegisterBrokenAuthHeaderProvider(auth.Url)
	configuration := clientcredentials.Config{
		ClientID:     auth.Client,
		ClientSecret: auth.Secret,
		TokenURL:     auth.Url,
	}

	return configuration, nil
}
