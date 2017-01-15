package cluster

import (
	"sync"

	"github.com/AsynkronIT/protoactor-go/process"
)

type pidCache struct {
	lock  sync.RWMutex
	Cache map[string]*process.ID
}

func (c *pidCache) Get(key string) *process.ID {
	c.lock.RLock()
	pid := c.Cache[key]
	c.lock.RUnlock()
	return pid
}

func (c *pidCache) Add(key string, pid *process.ID) {
	c.lock.Lock()
	c.Cache[key] = pid
	c.lock.Unlock()
}

var cache = &pidCache{
	Cache: make(map[string]*process.ID),
}
