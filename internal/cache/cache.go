package cache

import (
	"sync"
	"time"
)

type Item[T any] struct {
	Value     T
	ExpiresAt time.Time
}

type Cache[T any] struct {
	mu   sync.RWMutex
	data map[string]Item[T]
}

func New[T any]() *Cache[T] { return &Cache[T]{data: make(map[string]Item[T])} }

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	it, ok := c.data[key]
	if !ok || time.Now().After(it.ExpiresAt) {
		var zero T
		return zero, false
	}
	return it.Value, true
}

func (c *Cache[T]) Set(key string, val T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = Item[T]{Value: val, ExpiresAt: time.Now().Add(ttl)}
}
