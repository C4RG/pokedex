package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries map[string]cacheEntry
	mu      *sync.RWMutex
}

type cacheEntry struct {
	created time.Time
	val     []byte
}

func NewCache(i time.Duration) Cache {
	newCache := Cache{Entries: make(map[string]cacheEntry), mu: &sync.RWMutex{}}
	ticker := time.NewTicker(i)
	go func() {
		for range ticker.C {
			for k, v := range newCache.Entries {
				if time.Now().After(v.created.Add(i)) {
					newCache.mu.Lock()
					delete(newCache.Entries, k)
					newCache.mu.Unlock()
				}
			}
		}
	}()
	return newCache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	c.Entries[key] = cacheEntry{created: time.Now(), val: value}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	entry, exists := c.Entries[key]
	c.mu.RUnlock()
	if exists {
		return entry.val, exists
	}
	return []byte{}, exists
}
