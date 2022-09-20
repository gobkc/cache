package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	c := NewCache()
	for i := 0; i < 200; i++ {
		c.Set(fmt.Sprintf("a%v", i), i)
	}
	time.Sleep(5 * time.Second)
}
