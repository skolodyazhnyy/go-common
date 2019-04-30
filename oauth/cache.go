package oauth

import (
	"time"
)

type cache interface {
	// ShouldSet attempts to set value, in case value can not be set (for example due to an error)
	// this method should log error and return false
	ShouldSet(key string, value interface{}, ttl time.Duration) bool
	// ShouldGet attempts to get value, in case value can not be loaded (for example due to an error)
	// this method should log error and return false
	ShouldGet(key string, value interface{}) bool
}

type nopCache struct {
}

func (nopCache) ShouldSet(key string, value interface{}, ttl time.Duration) bool {
	return false
}

func (nopCache) ShouldGet(key string, value interface{}) bool {
	return false
}
