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
	aaa, _ := c.Load(5)
	fmt.Println(aaa)
	fmt.Println("Sleep 2")
	time.Sleep(2 * time.Second)
	fmt.Println("Get")
	aaa, _ = c.Load(5)
	fmt.Println(aaa)
	fmt.Println("Sleep 2")
	time.Sleep(2 * time.Second)
	fmt.Println("Get")
	aaa, _ = c.Load(5)
	fmt.Println(aaa)
	fmt.Println("Store 6")
	c.Set(6, "Hi")
	fmt.Println("Sleep 4")
	time.Sleep(4 * time.Second)
	fmt.Println("Store 7")
	c.Set(7, "Ho")
	fmt.Println("Get 6")
	aaa, _ = c.Load(6)
	fmt.Println(aaa)
}

func TestSet(t *testing.T) {
	c := NewCache(3*time.Second, false, 1*time.Second)
	c.Set(5, "Hello")
	c.Set(6, "Hello")
	c.Set(7, "Hello")
	c.Set(5, "Hello")
	c.Load(5)
}

func TestExtend(t *testing.T) {
	c := NewCache(3*time.Second, true, 0*time.Second)
	for i := 0; i < 10; i++ {
		c.Set(i, i)
	}
	time.Sleep(2 * time.Second)
	c.Load(5)
	c.Load(8)
	time.Sleep(2 * time.Second)
	c.Set(11, 11)
	aaa, _ := c.Load(5)
	fmt.Printf("Get 5:%v\n", aaa)
	aaa, _ = c.Load(6)
	fmt.Printf("Get 6:%v\n", aaa)
	aaa, _ = c.Load(8)
	fmt.Printf("Get 8:%v\n", aaa)
	time.Sleep(4 * time.Second)
	aaa, _ = c.Load(5)
	fmt.Printf("Get 5:%v\n", aaa)
	aaa, _ = c.Load(6)
	fmt.Printf("Get 6:%v\n", aaa)
	aaa, _ = c.Load(8)
	fmt.Printf("Get 8:%v\n", aaa)
	c.Set(11, 11)
}

func sleepUntil(t time.Time) {
	u := time.Until(t)
	if u < 0 {
		panic("Too late!")
	}
	time.Sleep(u)
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc: %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc: %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys: %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC: %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func RunRound(c *Cache, round int) {
	base := 1000000
	now := time.Now()
	for i := base * round; i < base*(round+1); i++ {
		c.Set(i, true)
	}

	fmt.Printf("Round: %v\t", round)
	PrintMemUsage()
	sleepUntil(now.Add(c.expiration).Add(-1000 * time.Millisecond))
}

func TestLarge(t *testing.T) {
	to := 5 * time.Second

	c := NewCache(to, true, 100*time.Microsecond)
	for r := 0; r < 30; r++ {
		RunRound(c, r)
	}
	c.Set(1, false)
	runtime.GC()
	aaa, _ := c.Load(1)
	fmt.Println(aaa)
	aaa, _ = c.Load(100)
	fmt.Println(aaa)
	fmt.Printf("Final: \t")
	PrintMemUsage()
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
