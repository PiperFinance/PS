package configs

import (
	"context"

	//"log"
	"os"
	"time"

	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/go-redis/redis/v8"
	"github.com/jellydator/ttlcache/v3"
	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

var (
	CacheManager *cache.ChainCache[any]
	TTLCache     = ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](15 * time.Second),
	)
)

func LoadCache() {
	REDIS_URL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		log.Error("Cache :: REDIS_URL env not found, defaulting to redis://127.0.0.1:6379")
		REDIS_URL = "redis://127.0.0.1:6379"
	}
	redisOption, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		log.Fatalf("Cache :: %s", err)
	}

	ctx := context.Background()
	gocacheStore := store.NewGoCache(gocache.New(5*time.Minute, 10*time.Minute))

	redisClient := redis.NewClient(redisOption)
	if err := redisClient.Ping(ctx); err != nil {
		log.Warn(err)
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
	_, err = CacheManager.Get(ctx, "Connected")
	if err != nil {
		log.Fatalf("unable to get cache key '%s' ", err)
	}
}
