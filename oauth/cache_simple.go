package oauth

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type simpleCache struct {
	lock          sync.RWMutex
	values        map[string]interface{}
	expire        map[string]time.Time
	lastSweep     time.Time
	sweepInterval time.Duration
}

// Set cache element
func (c *simpleCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.sweepInterval != 0 && time.Since(c.lastSweep) > c.sweepInterval {
		c.sweep()
	}

	c.values[key] = value

	if ttl != 0 {
		c.expire[key] = time.Now().Add(ttl)
	}

	return nil
}

// Get cache element
func (c *simpleCache) Get(key string, value interface{}) (bool, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if expires, ok := c.expire[key]; ok && expires.Before(time.Now()) {
		return false, nil
	}

	cached, ok := c.values[key]
	if !ok {
		return false, nil
	}

	ref := reflect.ValueOf(value)
	if ref.Kind() != reflect.Ptr {
		return false, errors.New("value receiver is not a pointer")
	}

	ref = ref.Elem()

	if !ref.CanSet() {
		return false, errors.New("value can not be set")
	}

	if !reflect.TypeOf(cached).AssignableTo(ref.Type()) {
		return false, errors.New("cache type does not match expected type")
	}

	ref.Set(reflect.ValueOf(cached))

	return true, nil
}

// sweep expired elements, this method assumes lock is set
func (c *simpleCache) sweep() {
	c.lastSweep = time.Now()

	for key, expire := range c.expire {
		if expire.Before(c.lastSweep) {
			delete(c.expire, key)
			delete(c.values, key)
		}
	}
}
