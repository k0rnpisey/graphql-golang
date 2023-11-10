package database

import (
	"github.com/go-gorm/caches/v2"
	"sync"
)

type Cacher struct {
	store *sync.Map
}

func (c *Cacher) init() {
	if c.store == nil {
		c.store = &sync.Map{}
	}
}

func (c *Cacher) Get(key string) *caches.Query {
	c.init()
	val, ok := c.store.Load(key)
	if !ok {
		return nil
	}

	return val.(*caches.Query)
}

func (c *Cacher) Store(key string, val *caches.Query) error {
	c.init()
	c.store.Store(key, val)
	return nil
}
