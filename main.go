package fixed_time_cache

import (
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
	if c.expiration <= 0 {
		return
	}
	c.ClearExpired()
	exptime := time.Now().Add(c.expiration)
	c.items.Store(key, Item{Object: val, Expiration: exptime})
	c.ExtendExpire(key, exptime)
}

func (c *Cache) Get(key interface{}) (val interface{}, ok bool) {
	if c.expiration <= 0 {
		return nil, false
	}
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
	if c.expiration <= 0 {
		return
	}
	c.timeouts_lock.Lock() //mode this item to the end of the linked list
	c.timeouts.Delete(key)
	c.timeouts.Set(key, exptime)
	c.timeouts_lock.Unlock()
}

func (c *Cache) ClearExpired() {
	if c.expiration <= 0 {
		return
	}
	if c.ClearCooldown != 0 && time.Now().Before(c.nextClear) {
		return
	}
	c.nextClear = time.Now().Add(c.ClearCooldown)
	need_clean := false

	c.timeouts_lock.RLock()
	pair := c.timeouts.Oldest()
	if pair != nil {
		if time.Now().After(pair.Value.(time.Time)) {
			need_clean = true
		}
	}
	c.timeouts_lock.RUnlock()
	if !need_clean {
		return
	}
	
	c.timeouts_lock.Lock()
	defer c.timeouts_lock.Unlock()

	for pair != nil {
		next := pair.Next()
		if time.Now().After(pair.Value.(time.Time)) {
			c.timeouts.Delete(pair.Key)
			c.items.Delete(pair.Key)
		} else {
			break
		}
		pair = next
	}

}
