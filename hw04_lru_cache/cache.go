package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

type Key string

type Cache interface {
	Set(Key, interface{}) bool
	Get(Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	dict     map[Key]*listItem
	mu       sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		dict:     make(map[Key]*listItem, capacity),
	}
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.dict[key]; ok {
		c.queue.MoveToFront(item)

		return item.Value.(*cacheItem).value, true
	}

	return nil, false
}

func (c *lruCache) Set(key Key, val interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.dict[key]; ok {
		item.Value.(*cacheItem).value = val
		c.queue.MoveToFront(item)

		return true
	}

	if len(c.dict) == c.capacity {
		tail := c.queue.Back()
		delete(c.dict, tail.Value.(*cacheItem).key)
		c.queue.Remove(tail)
	}

	item := &cacheItem{
		key:   key,
		value: val,
	}
	node := c.queue.PushFront(item)
	c.dict[key] = node

	return false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.dict = make(map[Key]*listItem, c.capacity)
}
