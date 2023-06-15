package multicaller

import (
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"

	"portfolio/configs"
	"portfolio/schema"
)

func GetChainsPairBalancesUnsafe(
	chainIds []schema.ChainId,
	wallet common.Address,
) map[schema.ChainId]schema.PairMapping {
	chunkedResultChannel := make(chan []ChunkCall[*big.Int])
	_res := make(map[schema.ChainId]schema.PairMapping)

	var totalChunkCount uint64
	// totalChunkCount = 0

	for _, chainId := range chainIds {
		_pairs := configs.ChainPairs(chainId)
		_multicall := configs.ChainMultiCall(chainId)

		if _multicall == nil || _pairs == nil {
			continue
		}
		atomic.AddUint64(
			&totalChunkCount,
			getPairBalances(PairBalanceCallOpt, chainId, *_multicall, _pairs, wallet, chunkedResultChannel, configs.ChainContextTimeOut(chainId)))
	}

	for chunkCalls := range chunkedResultChannel {
		// tmp := <-chunkedResultChannel
		if totalChunkCount > 0 {
			totalChunkCount--
		}

		balancePairResultParser(wallet, _res, chunkCalls)

		if totalChunkCount == 0 {
			break
		}
	}

	close(chunkedResultChannel)

	return _res
}
