package cache

import (
	"sync"
	"time"
)

var (
	once  sync.Once
	cache *Cache
)

func NewCache() ICache {
	once.Do(func() {
		cache = &Cache{
			expired:  5 * time.Second,
			interval: 1 * time.Second,
			data:     make(map[string]*Item),
		}
		go cache.GC()
	})
	return cache
}

func SetExpired(duration time.Duration) {
	if cache != nil {
		if duration == 0 {
			return
		}
		cache.expired = duration
	}
}

func SetInterval(duration time.Duration) {
	if cache != nil {
		if duration == 0 {
			return
		}
		cache.interval = duration
	}
}
