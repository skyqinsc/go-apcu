package apcu

import (
	"sync"
	"time"
)

const (
	//DefaultExpiration 默认失效时间10分钟
	DefaultExpiration time.Duration = 10 * time.Minute
)

var (
	GlobalCache *Cache
)

//Cache ...
type Cache struct {
	*cache
}

type cache struct {
	items             map[string]*Item
	mu                sync.RWMutex
	defaultExpiration time.Duration
	janitor           *janitorQueue
}

//NewCache 新建本地缓存
func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]*Item)
	c := &cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		janitor:           newJanitor(cleanupInterval, MaxWatcherSize),
	}
	c.janitor.run(c)
	C := &Cache{c}
	return C
}

func (c *cache) Store(key string, val interface{}, ttl time.Duration) {
	var expiredTime int64
	if ttl > 0 {
		expiredTime = time.Now().Add(ttl).UnixNano()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &Item{
		Object:     val,
		Expiration: expiredTime,
	}
	c.janitor.notify(key, expiredTime)
}

func (c *cache) Fetch(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if item.Expired() {
		return nil, false
	}

	return item.Object, true
}
