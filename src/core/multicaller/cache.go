package multicaller

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/jellydator/ttlcache/v3"
	"math/big"
	"portfolio/configs"
	"portfolio/schema"
)

type ChunkCallCacheKey struct {
	Wallet common.Address
	schema.Id
}
type ChunkCallsCacheKey struct {
	Wallet common.Address
	schema.ChainId
	chunkIndex int
	what       string
}

var ChunkCallsCache = ttlcache.New[ChunkCallsCacheKey, []ChunkCall[*big.Int]](
	ttlcache.WithTTL[ChunkCallsCacheKey, []ChunkCall[*big.Int]](configs.ChunkCallCacheTTL),
)
