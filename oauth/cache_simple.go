package oauth

import (
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

// ShouldSet cache element
func (c *simpleCache) ShouldSet(key string, value interface{}, ttl time.Duration) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.sweepInterval != 0 && time.Since(c.lastSweep) > c.sweepInterval {
		c.sweep()
	}

	c.values[key] = value

	if ttl != 0 {
		c.expire[key] = time.Now().Add(ttl)
	}

	return true
}

// ShouldGet cache element
func (c *simpleCache) ShouldGet(key string, value interface{}) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if expires, ok := c.expire[key]; ok && expires.Before(time.Now()) {
		return false
	}

	cached, ok := c.values[key]
	if !ok {
		return false
	}

	ref := reflect.ValueOf(value)
	if ref.Kind() != reflect.Ptr {
		return false
	}

	ref = ref.Elem()

	if !ref.CanSet() {
		return false
	}

	if !reflect.TypeOf(cached).AssignableTo(ref.Type()) {
		return false
	}

	ref.Set(reflect.ValueOf(cached))

	return true
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
