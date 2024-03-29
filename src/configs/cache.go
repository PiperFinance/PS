package configs

import (
	"context"
	"fmt"
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/go-redis/redis/v8"
	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	//"log"
	"sync"

	"os"
	"time"
)

var (
	//CacheManager *cache.Cache[any]
	onceForGoCache sync.Once
	CacheManager   *cache.ChainCache[any]

	GoCacheCleanupInterval = 10 * time.Minute
	GoCacheDefaultTTL      = 5 * time.Minute
	ChunkCallCacheTTL      = 1 * time.Minute
)

func init() {

	onceForGoCache.Do(func() {

		RedisUrl, ok := os.LookupEnv("REDIS_URL")
		if !ok {
			log.Error("Cache :: REDIS_URL env not found, defaulting to redis://127.0.0.1:6379")
			RedisUrl = "redis://127.0.0.1:6379"
		}
		redisOption, err := redis.ParseURL(RedisUrl)
		if err != nil {
			log.Fatalf("Cache :: %s", err)
		}

		ctx := context.Background()
		gocacheStore := store.NewGoCache(gocache.New(GoCacheDefaultTTL, GoCacheCleanupInterval))

		redisClient := redis.NewClient(redisOption)
		if err := redisClient.Ping(ctx); err != nil {
			log.Error(err)
			CacheManager = cache.NewChain[any](
				cache.New[any](gocacheStore),
			)
		} else {
			redisStore := store.NewRedis(redisClient, nil)
			CacheManager = cache.NewChain[any](
				cache.New[any](gocacheStore),
				cache.New[any](redisStore),
			)
		}

		err = CacheManager.Set(ctx, "Connected", "YES", store.WithExpiration(15*time.Second))

		if err != nil {
			panic(err)
		}
		value, err := CacheManager.Get(ctx, "Connected")
		if err != nil {
			log.Fatalf("unable to get cache key '%s' ", err)
		}
		fmt.Printf("%#+v\n", value)
	})
}
