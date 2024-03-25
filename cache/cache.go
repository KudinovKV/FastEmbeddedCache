package cache

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("element not found")

type Driver struct {
	data map[string]*item
	mu   sync.RWMutex

	defaultTTL time.Duration
}

func (d *Driver) Set(key string, value any, ttl time.Duration) {
	if ttl == 0 {
		ttl = d.defaultTTL
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.data[key] = newItem(value, ttl)
}

func (d *Driver) Get(key string) (any, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	cacheItem, ok := d.data[key]
	if !ok {
		return nil, ErrNotFound
	}

	if cacheItem.isExpired() {
		go d.Delete(key)

		return nil, ErrNotFound
	}

	return cacheItem.value, nil
}

func (d *Driver) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.data, key)
}

func New(defaultTTL time.Duration) *Driver {
	return &Driver{
		data:       make(map[string]*item),
		defaultTTL: defaultTTL,
	}
}
