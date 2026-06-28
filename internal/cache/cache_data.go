package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type CacheData struct {
	data  map[string]cacheEntry
	mutex sync.RWMutex
}

func NewCacheData() *CacheData {
	return &CacheData{
		data: make(map[string]cacheEntry),
	}
}

func (c *CacheData) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}

func (c *CacheData) Set(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *CacheData) reapLoop() {
	//time.Ticker to periodically clean up old entries
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.reap()
	}
}

func (c *CacheData) reap() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.Sub(entry.createdAt) > 5*time.Second {
			delete(c.data, key)
		}
	}
}
