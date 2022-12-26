package cache

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type CacheMetadata struct {
	ExpirationTime time.Duration
}

type CacheConfig struct {
	metadata CacheMetadata
	once     sync.Once
}

// Private
func (cacheConfig *CacheConfig) lazyInit() {
	cacheConfig.once.Do(func() {
		numberOfExpirationTimeSeconds, err := strconv.Atoi(os.Getenv("CACHE_EXPIRATION"))
		if err != nil {
			panic(err)
		}
		expirationTime := time.Duration(numberOfExpirationTimeSeconds) * time.Second

		cacheConfig.metadata.ExpirationTime = expirationTime
	})
}

// Public
func (cacheConfig *CacheConfig) GetMetadata() CacheMetadata {
	cacheConfig.lazyInit()
	return cacheConfig.metadata
}

var Config = &CacheConfig{}
