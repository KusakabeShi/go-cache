package fixed_time_cache

import (
	"container/list"
	"fmt"
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
	timeouts   *list.List
	mu         sync.RWMutex
}

func NewCache(defaultExpiration time.Duration) *Cache {
	return &Cache{
		expiration: defaultExpiration,
		items:      map[interface{}]Item{},
		timeouts:   list.New(),
	}
}

func (c *Cache) Set(key interface{}, val interface{}) {
	c.ClearExpired()
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

func (c *Cache) Get(key interface{}) (val interface{}, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.items[key]; ok {
		if val.Expiration.After(time.Now()) {
			return val.Object, true
		}
	}
	return nil, false
}

func (c *Cache) ClearExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for c.timeouts.Len() > 0 {
		first := c.timeouts.Front()
		switch first_elem := first.Value.(type) {
		case timeout:
			if time.Now().After(first_elem.Expiration) {
				delete(c.items, first_elem.Key)
				c.timeouts.Remove(first)
			} else {
				return
			}
		}
	}
}

func fixed_time_cache_test() {

	c := NewCache(3 * time.Second)
	c.Set(5, "Hello")
	aaa, ok := c.Get(5)
	switch a := aaa.(type) {
	case string:
		fmt.Println(a, ok)
	}

	aaa, ok = c.Get(5)
	switch a := aaa.(type) {
	case string:
		fmt.Println(a, ok)
	}
	return
}
