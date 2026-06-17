package registry

import (
	"realtime_blockchain/internal/model"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	store map[string]model.PriceEntry
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]model.PriceEntry),
	}
}

func (c *Cache) Get(addr string) (model.PriceEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.store[addr]
	if !ok {
		return model.PriceEntry{}, false
	}

	return v, true
}

func (c *Cache) Set(addr string, p model.PriceEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[addr] = p
}
