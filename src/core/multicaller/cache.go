package multicaller

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/jellydator/ttlcache/v3"
	"math/big"
	"portfolio/schema"
	"time"
)

type ChunkCallCacheKey struct {
	Wallet common.Address
	schema.Id
}
type ChunkCallsCacheKey struct {
	Wallet common.Address
	schema.ChainId
	chunkIndex int
}

var (
	ChunkCallCacheTTL = 15 * time.Second

	//ChunkCallsCache
	// True For OK , False For Failed
	ChunkCallsCache = ttlcache.New[ChunkCallsCacheKey, []ChunkCall[*big.Int]](
		ttlcache.WithTTL[ChunkCallsCacheKey, []ChunkCall[*big.Int]](ChunkCallCacheTTL),
	)
)

//func ChunkCallsCache[T comparable]() {
//	switch T.(type) {
//	case big.Int:
//
//	}
//	ttlcache.New[ChunkCallsCacheKey, []ChunkCall[T]](
//		ttlcache.WithTTL[ChunkCallsCacheKey, []ChunkCall[T]](ChunkCallCacheTTL),
//	)
//}
