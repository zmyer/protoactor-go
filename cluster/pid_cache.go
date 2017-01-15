package cluster

import (
	"sync"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type pidCache struct {
	lock  sync.RWMutex
	Cache map[string]*process.PID
}

func (c *pidCache) Get(key string) *process.PID {
	c.lock.RLock()
	pid := c.Cache[key]
	c.lock.RUnlock()
	return pid
}

func (c *pidCache) Add(key string, pid *process.PID) {
	c.lock.Lock()
	c.Cache[key] = pid
	c.lock.Unlock()
}

var cache = &pidCache{
	Cache: make(map[string]*process.PID),
}
