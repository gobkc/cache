# cache
a generic and lightweight cache component

### Example
````
package main

import (
	"fmt"
	memCache "github.com/gobkc/cache"
)

func main() {
	c := memCache.NewCache()
	c.Set("aaa", 123)
	loadVal := 0
	if err := c.GetObject("aaa", &loadVal); err != nil {
		fmt.Println(err)
	}
	fmt.Println("load val:", loadVal)
}
````