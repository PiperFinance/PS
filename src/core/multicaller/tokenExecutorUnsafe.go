package multicaller

import (
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"

	"portfolio/configs"
	"portfolio/schema"
)

func GetChainsTokenBalancesUnsafe(
	chainIds []schema.ChainId,
	wallet common.Address,
) (map[schema.ChainId]schema.TokenMapping, error) {
	chunkedResultChannel := make(chan []ChunkCall[*big.Int])
	_res := make(map[schema.ChainId]schema.TokenMapping)

	var totalChunkCount uint64
	// totalChunkCount = 0

	for _, chainId := range chainIds {
		_tokens := configs.ChainTokens(chainId)
		_multicall, err := configs.ChainMultiCall(chainId)
		if err != nil {
			return nil, err
		}

		if _multicall == nil || _tokens == nil {
			continue
		}
		atomic.AddUint64(
			&totalChunkCount,
			getTokenBalances(TokenBalanceCallOpt, chainId, _multicall, _tokens, wallet, chunkedResultChannel, configs.ChainContextTimeOut(chainId)))
	}
	if totalChunkCount == 0 {
		return _res, nil
	}
	for chunkCalls := range chunkedResultChannel {
		if totalChunkCount > 0 {
			totalChunkCount--
		}
		balanceTokenResultParser(wallet, _res, chunkCalls)
		if totalChunkCount == 0 {
			break
		}
	}

	close(chunkedResultChannel)

	return _res, nil
}
