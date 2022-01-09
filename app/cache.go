package app

import (
	"sync"
)

// Cache 并发控制
type Cache struct {
	mu         sync.Mutex
	lru        *lru
	cacheBytes int64
}

func NewCache(cacheBytes int64, OnEvicted func(key string, value Value)) *Cache {
	return &Cache{
		cacheBytes: cacheBytes,
		lru:        New(cacheBytes, OnEvicted),
	}
}

// Set 添加缓存
func (c *Cache) Set(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = New(c.cacheBytes, nil)
	}

	c.lru.Add(key, value)
}

// Get 获取缓存
func (c *Cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value, ok := c.lru.Get(key); ok {
		return value.(ByteView), true
	}

	return
}
