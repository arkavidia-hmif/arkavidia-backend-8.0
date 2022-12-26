package cache

import (
	"sync"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"

	cacheConfig "arkavidia-backend-8.0/competition/config/cache"
)

// TODO: Tambahkan cache layer pada route dengan method GET
// REFERENCE: https://github.com/gin-contrib/cache
// ASSIGNED TO: @patrickamadeus

type LocalCache struct {
	store *persistence.InMemoryStore
	once  sync.Once
}

// Private
func (localCache *LocalCache) lazyInit() {
	localCache.once.Do(func() {
		config := cacheConfig.Config.GetMetadata()
		localCache.store = persistence.NewInMemoryStore(config.ExpirationTime)
	})
}

// Public
func (localCache *LocalCache) GetHandlerFunc(handle gin.HandlerFunc) gin.HandlerFunc {
	localCache.lazyInit()
	config := cacheConfig.Config.GetMetadata()
	return cache.CachePage(localCache.store, config.ExpirationTime, handle)
}

var Store = &LocalCache{}
