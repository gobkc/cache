package cache

import "time"

type Item struct {
	v      interface{}
	create time.Time
	out    time.Time
}

func (i *Item) IsExpired() bool {
	return time.Now().Unix() > i.out.Unix()
}
