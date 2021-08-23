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
	Key    interface{}
	Expire time.Time
}

type Cache struct {
	expiration    time.Duration
	items         sync.Map
	timeouts      *list.List
	nextClear     time.Time
	ClearCooldown time.Duration
}

func NewCache(defaultExpiration time.Duration, clearcooldown time.Duration) *Cache {
	return &Cache{
		expiration:    defaultExpiration,
		timeouts:      list.New(),
		ClearCooldown: clearcooldown,
	}
}

func (c *Cache) Set(key interface{}, val interface{}) {
	c.ClearExpired()
	exptime := time.Now().Add(c.expiration)
	c.items.Store(key, Item{Object: val, Expiration: exptime})
	c.timeouts.PushBack(timeout{Key: key, Expire: exptime})
}

func (c *Cache) Get(key interface{}) (val interface{}, ok bool) {
	if val, ok := c.items.Load(key); ok {
		if val.(Item).Expiration.After(time.Now()) {
			return val.(Item).Object, true
		} else {
			c.items.Delete(key)
		}
	}
	return nil, false
}

func (c *Cache) ClearExpired() {
	if c.ClearCooldown != 0 && time.Now().Before(c.nextClear) {
		return
	}
	c.nextClear = time.Now().Add(c.ClearCooldown)
	for c.timeouts.Len() > 0 {
		first := c.timeouts.Front()
		switch first_elem := first.Value.(type) {
		case timeout:
			if time.Now().After(first_elem.Expire) {
				c.items.Delete(first_elem.Key)
				c.timeouts.Remove(first)
			} else {
				return
			}
		}
	}
}

func example() {

	c := NewCache(3*time.Second, 1*time.Second)
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
