package cache

import (
	"fmt"
	"testing"
)

func TestNewCache(t *testing.T) {
	c := NewCache()
	for i := 0; i < 1000000; i++ {
		go c.Set(fmt.Sprintf("a%v", i), i)
	}
}
