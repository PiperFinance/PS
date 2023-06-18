package multicaller

import (
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/core/utils"
	"portfolio/schema"
)

var TokenBalanceCallOpt ChunkedCallOpts

func init() {
	TokenBalanceCallOpt = ChunkedCallOpts{W3CallOpt: nil, ChunkSize: 10}
}

// getTokenBalances Wallet balance based on given token ( Faster if chunks is used)
// Does not sort + only respond with tokens with balance
func getTokenBalances(
	callOpts ChunkedCallOpts,
	id schema.ChainId,
	multiCaller *Multicall.MulticallCaller,
	tokens schema.TokenMapping,
	wallet common.Address,
	chunkedResultChannel chan []ChunkCall[*big.Int],
	callTimeout time.Duration,
) uint64 {
	allCalls := genTokenBalanceCalls(tokens, wallet)
	chunkedCalls := utils.Chunks[ChunkCall[*big.Int]](allCalls, callOpts.ChunkSize)

	for i, indexCalls := range chunkedCalls {
		cachedChunkCalls := ChunkCallsCache.Get(ChunkCallsCacheKey{wallet, id, i})
		if cachedChunkCalls != nil && !cachedChunkCalls.IsExpired() {
			go func() {
				chunkedResultChannel <- cachedChunkCalls.Value()
			}()
		} else {
			go execute(i, id, wallet, multiCaller, indexCalls, chunkedResultChannel, callTimeout)
		}
	}
	return uint64(len(chunkedCalls))
}

func balanceTokenResultParser(wallet common.Address, result map[schema.ChainId]schema.TokenMapping, chunk []ChunkCall[*big.Int]) {
	for _, call := range chunk {

		// TODO In case error occurred at rpc level
		if call.Err != nil {
		}

		if !call.CallRes.Success || call.ParsedCallRes == nil {
			continue
		}
		chainId := call.ChainId
		if result[chainId] == nil {
			result[chainId] = make(schema.TokenMapping)
		}
		_tokenId := schema.TokenId(call.Id)
		_token := configs.ChainTokens(chainId)[_tokenId].Copy()

		// Token Balance
		utils.MustParseBal(call.ParsedCallRes, _token)

		result[chainId][_tokenId] = *_token
	}
	_ = wallet
}

func GetChainsTokenBalances(
	chainIds []schema.ChainId,
	wallet common.Address,
	callTimeout time.Duration,
) (map[schema.ChainId]schema.TokenMapping, error) {
	chunkedResultChannel := make(chan []ChunkCall[*big.Int])
	_res := make(map[schema.ChainId]schema.TokenMapping)

	var totalChunkCount uint64

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
			getTokenBalances(TokenBalanceCallOpt, chainId, _multicall, _tokens, wallet, chunkedResultChannel, callTimeout))
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
