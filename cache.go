package cache

import (
	"container/list"
	"sync"
	"time"
)

type Item struct {
	Object     interface{}
	Expiration time.Time
}

type timeout struct {
	Key        interface{}
	Expiration time.Time
}

type Cache struct {
	expiration time.Duration
	items      map[interface{}]Item
	timeouts   list.List
	mu         sync.RWMutex
}

func (*Cache) NewCache(defaultExpiration time.Duration) *Cache {
	return &Cache{
		expiration: defaultExpiration,
		items:      map[interface{}]Item{},
		timeouts:   *list.New().Init(),
	}
}

func (c *Cache) Set(key interface{}, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	exptime := time.Now().Add(c.expiration)
	c.items[key] = Item{
		Object:     val,
		Expiration: exptime,
	}
	c.timeouts.PushBack(timeout{
		Key:        key,
		Expiration: exptime,
	})
}

func (c *Cache) Get(key interface{}) interface{} {
	c.clearExpired()
	c.mu.RLock()
	defer c.mu.Unlock()
	return c.items[key]
}

func (c *Cache) clearExpired() {
	for c.timeouts.Len() > 0 {
		first := c.timeouts.Front()
		switch first_elem := first.Value.(type) {
		case timeout:
			if time.Now().After(first_elem.Expiration) {
				delete(c.items, first_elem.Key)
				c.timeouts.Remove(first)
			} else {
				break
			}
		}
	}
}
