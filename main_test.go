package fixed_time_cache

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestExpire(t *testing.T) {
	c := NewCache(3*time.Second, false, 1*time.Second)
	fmt.Println("Store")
	c.Set(5, "Hello")
	fmt.Println("Get")
	aaa, ok := c.Get(5)
	switch a := aaa.(type) {
	case string:
		fmt.Println(a, ok)
	case nil:
		fmt.Println(a, ok)
	}
	fmt.Println("Sleep 2")
	time.Sleep(2 * time.Second)
	fmt.Println("Get")
	aaa, ok = c.Get(5)
	switch a := aaa.(type) {
	case string:
		fmt.Println(a, ok)
	case nil:
		fmt.Println(a, ok)
	}
	fmt.Println("Sleep 2")
	time.Sleep(2 * time.Second)
	fmt.Println("Get")
	aaa, ok = c.Get(5)
	switch a := aaa.(type) {
	case string:
		fmt.Println(a, ok)
	case nil:
		fmt.Println(a, ok)
	}
	fmt.Println("Store 6")
	c.Set(6, "Hi")
	fmt.Println("Sleep 4")
	time.Sleep(4 * time.Second)
	fmt.Println("Store 7")
	c.Set(7, "Ho")
	fmt.Println("Get 6")
	aaa, ok = c.Get(6)
	switch a := aaa.(type) {
	case string:
		fmt.Println(a, ok)
	case nil:
		fmt.Println(a, ok)
	}
}

func TestExtend(t *testing.T) {
	c := NewCache(1*time.Second, true, 0*time.Second)
	for i := 0; i < 10000; i++ {
		c.Set(i, "Hello")
		c.Set(i, "Hi")
		c.Set(i, "Hey")
	}
	for i := 0; i < 10000; i++ {
		c.Get(i)
		c.Get(i)
		c.Get(i)
	}
	time.Sleep(2 * time.Second)
	for i := 0; i < 10000; i++ {
		c.Set(5, "Hello")
		c.Set(8, "Ho")
	}
	for i := 0; i < 10000; i++ {
		c.Get(5)
		c.Get(8)
	}
	aaa, _ := c.Get(5)
	fmt.Printf("Get 5:%v\n", aaa)
	aaa, _ = c.Get(6)
	fmt.Printf("Get 6:%v\n", aaa)
	aaa, _ = c.Get(5)
	fmt.Printf("Get 5:%v\n", aaa)
}

func sleepUntil(t time.Time) {
	u := time.Until(t)
	if u < 0 {
		panic("Too late!")
	}
	time.Sleep(u)
}

func TestLarge(t *testing.T) {
	to := 5 * time.Second

	c := NewCache(to, true, 1*time.Second)
	base := 1000000
	n := time.Now()
	for i := 0; i < base; i++ {
		c.Set(i, true)
	}
	sleepUntil(n.Add(to).Add(time.Second))
	n = time.Now()
	for i := base; i < 2*base; i++ {
		c.Set(i, true)
	}
	sleepUntil(n.Add(to).Add(time.Second))
	n = time.Now()
	for i := 2 * base; i < 3*base; i++ {
		c.Set(i, true)
	}
	sleepUntil(n.Add(to).Add(time.Second))
	n = time.Now()
	c.ClearExpired()
	for i := 2 * base; i < 3*base; i++ {
		c.Set(i, true)
	}
	sleepUntil(n.Add(to).Add(time.Second))
	n = time.Now()
	c.ClearExpired()
	for i := 2 * base; i < 3*base; i++ {
		c.Set(i, true)
	}
	sleepUntil(n.Add(to).Add(time.Second))
	n = time.Now()
	c.ClearExpired()
	for i := 2 * base; i < 3*base; i++ {
		c.Set(i, true)
	}
	sleepUntil(n.Add(to).Add(time.Second))
	n = time.Now()
	c.ClearExpired()
	for i := 2 * base; i < 3*base; i++ {
		c.Set(i, true)
	}
	sleepUntil(n.Add(to).Add(time.Second))
	n = time.Now()
	c.ClearExpired()
	for i := 2 * base; i < 3*base; i++ {
		c.Set(i, true)
	}
	sleepUntil(n.Add(to).Add(time.Second))
	n = time.Now()
	c.ClearExpired()

	aaa, _ := c.Get(1)
	fmt.Println(aaa)
	aaa, _ = c.Get(100)
	fmt.Println(aaa)
}

func TestSynaMap(t *testing.T) {
	var s sync.Map
	for i := 0; i < 10000000; i++ {
		s.Store(i, true)
	}
	for i := 0; i < 10000000; i++ {
		s.Delete(i)
	}
	for i := 10000000; i < 20000000; i++ {
		s.Store(i, true)
	}
	for i := 10000000; i < 20000000; i++ {
		s.Delete(i)
	}
	for i := 20000000; i < 30000000; i++ {
		s.Store(i, true)
	}
	for i := 20000000; i < 30000000; i++ {
		s.Delete(i)
	}

	runtime.GC()
	aaa, _ := s.Load(1)
	fmt.Println(aaa)
	aaa, _ = s.Load(100)
	fmt.Println(aaa)
}
