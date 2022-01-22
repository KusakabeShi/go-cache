package fixed_time_cache

import (
	"fmt"
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
		c.Set(5, "Hello")
		c.Set(6, "Hi")
		c.Set(7, "Hey")
	}
	for i := 0; i < 10000; i++ {
		c.Get(5)
		c.Get(6)
		c.Get(7)
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
