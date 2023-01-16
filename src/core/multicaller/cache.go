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

var (
	ChunkCallCacheTTL = 15 * time.Second
	ChunkCallCache    = ttlcache.New[ChunkCallCacheKey, ChunkCall[*big.Int]](
		ttlcache.WithTTL[ChunkCallCacheKey, ChunkCall[*big.Int]](ChunkCallCacheTTL),
	)
	FailedChunkCallCache = ttlcache.New[ChunkCallCacheKey, ChunkCall[*big.Int]](
		ttlcache.WithTTL[ChunkCallCacheKey, ChunkCall[*big.Int]](ChunkCallCacheTTL),
	)
)
