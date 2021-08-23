package fixed_time_cache

import (
	"fmt"
	"sync"
	"time"

	orderedmap "github.com/wk8/go-ordered-map"
)

type Item struct {
	Object     interface{}
	Expiration time.Time
}

type Cache struct {
	expiration    time.Duration
	ExtendOnGet   bool
	items         sync.Map
	timeouts      *orderedmap.OrderedMap
	timeouts_lock sync.RWMutex
	nextClear     time.Time
	ClearCooldown time.Duration
}

func NewCache(defaultExpiration time.Duration, extendOnGet bool, clearcooldown time.Duration) *Cache {
	return &Cache{
		expiration:    defaultExpiration,
		ExtendOnGet:   extendOnGet,
		timeouts:      orderedmap.New(),
		ClearCooldown: clearcooldown,
	}
}

func (c *Cache) Set(key interface{}, val interface{}) {
	c.ClearExpired()
	exptime := time.Now().Add(c.expiration)
	c.items.Store(key, Item{Object: val, Expiration: exptime})
	c.ExtendExpire(key, exptime)
}

func (c *Cache) Get(key interface{}) (val interface{}, ok bool) {
	if val, ok := c.items.Load(key); ok {
		if val.(Item).Expiration.After(time.Now()) {
			if c.ExtendOnGet {
				c.ExtendExpire(key, time.Now().Add(c.expiration))
			}
			return val.(Item).Object, true
		} else {
			c.items.Delete(key)
		}
	}
	return nil, false
}

func (c *Cache) ExtendExpire(key interface{}, exptime time.Time) {
	c.timeouts_lock.Lock() //mode this item to the end of the linked list
	c.timeouts.Delete(key)
	c.timeouts.Set(key, exptime)
	c.timeouts_lock.Unlock()
}

func (c *Cache) ClearExpired() {
	if c.ClearCooldown != 0 && time.Now().Before(c.nextClear) {
		return
	}
	c.nextClear = time.Now().Add(c.ClearCooldown)
	c.timeouts_lock.RLock()
	for pair := c.timeouts.Oldest(); pair != nil; pair = pair.Next() {
		if time.Now().After(pair.Value.(time.Time)) {
			c.timeouts_lock.RUnlock()
			c.timeouts_lock.Lock()
			c.timeouts.Delete(pair.Key)
			c.items.Delete(pair.Key)
			c.timeouts_lock.Unlock()
			c.timeouts_lock.RLock()
		}
	}
	c.timeouts_lock.RUnlock()
}

func example() {

	c := NewCache(3*time.Second, false, 1*time.Second)
	c.Set(5, "Hello")
	aaa, ok := c.Get(5)
	switch a := aaa.(type) {
	case string:
		fmt.Println(a, ok)
	case nil:
		fmt.Println(a, ok)
	}

	time.Sleep(2 * time.Second)
	aaa, ok = c.Get(5)
	switch a := aaa.(type) {
	case string:
		fmt.Println(a, ok)
	case nil:
		fmt.Println(a, ok)
	}

	time.Sleep(2 * time.Second)
	aaa, ok = c.Get(5)
	switch a := aaa.(type) {
	case string:
		fmt.Println(a, ok)
	case nil:
		fmt.Println(a, ok)
	}

	c.Set(6, "Hi")
	time.Sleep(4 * time.Second)
	c.Set(7, "Ho")
	aaa, ok = c.Get(6)
	switch a := aaa.(type) {
	case string:
		fmt.Println(a, ok)
	case nil:
		fmt.Println(a, ok)
	}
	return
}
