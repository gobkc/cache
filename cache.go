package cache

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"sync"
	"time"
)

type Cache struct {
	lock     sync.RWMutex
	expired  time.Duration
	interval time.Duration
	data     map[string]*Item
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
	return nil, nil
}

func (c *Cache) GetObject(key string, data interface{}) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if v, findOk := c.data[key]; findOk {
		if v.IsExpired() {
			return errors.New("val is expired")
		}
		switch reflect.TypeOf(data).Kind() {
		case reflect.Ptr:
			switch reflect.ValueOf(data).Elem().Kind() {
			case reflect.Struct:
				reflect.ValueOf(data).Elem().Set(reflect.ValueOf(v.v))
				return nil
			default:
				return errors.New("dest must be a struct pointer")
			}
		default:
			return errors.New("dest must be a struct pointer")
		}
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
			if err := c.Delete(key); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Cache) GC() {
	for {
		data, _ := c.GetALL()
		b, _ := json.MarshalIndent(data, "", "\t")
		log.Println("gc data:", string(b))
		go func() {
			if err := c.DeleteExpired(); err != nil {
				log.Panicln(err.Error())
			}
		}()
		time.Sleep(c.interval)
	}
}
