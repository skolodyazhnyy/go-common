package oauth

import (
	"reflect"
	"testing"
	"time"
)

func TestSimpleCache(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value []string
		ttl   time.Duration
	}{
		{name: "simple string", key: "string", value: []string{"foo"}},
	}

	cache := &simpleCache{
		values:        make(map[string]interface{}),
		expire:        make(map[string]time.Time),
		lastSweep:     time.Now(),
		sweepInterval: time.Second,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := cache.Set(test.key, test.value, test.ttl); err != nil {
				t.Error("An error occurred while setting cache key:", err)
			}

			var got []string

			ok, err := cache.Get(test.key, &got)
			if err != nil {
				t.Fatal("Value can not be read:", err.Error())
			}

			if !ok {
				t.Fatal("Key is not present in the cache")
			}

			if !reflect.DeepEqual(got, test.value) {
				t.Errorf("Value from cache does not match value in the cache: want %v, got %v", test.value, got)
			}
		})
	}
}
