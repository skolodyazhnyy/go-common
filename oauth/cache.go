package oauth

import (
	"time"
)

type cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string, value interface{}) (bool, error)
}

type nopCache struct {
}

func (nopCache) Set(key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (nopCache) Get(key string, value interface{}) (bool, error) {
	return false, nil
}
