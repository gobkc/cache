package cache

import (
	"errors"
	"log"
	"reflect"
	"sync"
	"time"
)

type Cache struct {
	lock        sync.RWMutex
	expired     time.Duration
	interval    time.Duration
	concurrency int64
	data        map[string]*Item
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if v, findOk := c.data[key]; findOk {
		if v.IsExpired() {
			return nil, errors.New("val is expired")
		}
		return v.v, nil
	}
	return nil, nil
}

func (c *Cache) GetALL() (map[string]interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	data := make(map[string]interface{})
	for key, item := range c.data {
		data[key] = item
	}
	return data, nil
}

func (c *Cache) GetObject(key string, data interface{}) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if v, findOk := c.data[key]; findOk {
		if v.IsExpired() {
			return errors.New("val is expired")
		}
		destKind := reflect.TypeOf(data).Kind()
		if destKind != reflect.Ptr {
			return errors.New("dest must be a pointer")
		}
		vKind := reflect.TypeOf(v.v).Kind()
		if vKind == reflect.Ptr {
			vKind = reflect.TypeOf(v.v).Elem().Kind()
		}
		if reflect.TypeOf(data).Elem().Kind() != vKind {
			return errors.New("val type error")
		}
		reflect.ValueOf(data).Elem().Set(reflect.ValueOf(v.v))
		return nil
	}
	return nil
}

func (c *Cache) Set(key string, value interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if key == "" {
		return errors.New("key is empty")
	}
	c.data[key] = &Item{
		v:      value,
		create: time.Now(),
		out:    time.Now().Add(c.expired),
	}
	return nil
}

func (c *Cache) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, findOk := c.data[key]; !findOk {
		return errors.New("key not exist")
	}
	delete(c.data, key)
	return nil
}

func (c *Cache) DeleteExpired() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	for key, item := range c.data {
		if item.IsExpired() {
			if _, findOk := c.data[key]; !findOk {
				return errors.New("key not exist")
			}
			delete(c.data, key)
		}
	}
	return nil
}

func (c *Cache) GC() {
	limit := make(chan struct{}, c.concurrency)
	for {
		limit <- struct{}{}
		go c.toGC(&limit)
		time.Sleep(c.interval)
	}
}

func (c *Cache) toGC(limit *chan struct{}) {
	if err := c.DeleteExpired(); err != nil {
		log.Println(err.Error())
	}
	<-*limit
}
