package multicaller

import (
	"math/big"
	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/core/utils"
	"portfolio/schema"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

var (
	TokenBalanceCallOpt ChunkedCallOpts
)

func init() {
	TokenBalanceCallOpt = ChunkedCallOpts{W3CallOpt: nil, ChunkSize: 5}
}

// getTokenBalances Wallet balance based on given token ( Faster if chunks is used)
// Does not sort + only respond with tokens with balance
func getTokenBalances(
	callOpts ChunkedCallOpts,
	id schema.ChainId,
	multiCaller Multicall.MulticallCaller,
	tokens schema.TokenMapping,
	wallet common.Address,
	chunkedResultChannel chan []ChunkCall[*big.Int]) uint64 {

	allCalls := genTokenBalanceCalls(tokens, wallet)
	chunkedCalls := utils.Chunks[ChunkCall[*big.Int]](allCalls, callOpts.ChunkSize)

	for i, indexCalls := range chunkedCalls {
		cachedChunkCalls := ChunkCallsCache.Get(ChunkCallsCacheKey{wallet, id, i, "t"})
		if cachedChunkCalls != nil && !cachedChunkCalls.IsExpired() {
			go func() {
				chunkedResultChannel <- cachedChunkCalls.Value()
			}()
		} else {
			go execute("t", i, id, wallet, multiCaller, indexCalls, chunkedResultChannel)
		}
	}

	return uint64(len(chunkedCalls))
}

func balanceTokenResultParser(wallet common.Address, result map[schema.ChainId]schema.TokenMapping, chunk []ChunkCall[*big.Int]) {
	for _, call := range chunk {

		// In case error occurred at rpc level
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
		decimal := configs.DecimalPowTen(_token.Detail.Decimals)
		b := new(big.Float).SetInt(call.ParsedCallRes)
		b = b.Quo(b, new(big.Float).SetInt(decimal))
		if b.IsInf() {
			log.Errorf("[INF] @ (%d,%s) ::: cnResp  %s ", _token.Detail.Decimals, _token.Detail.Address, call.ParsedCallRes)
		}
		_token.Balance = *b
		_token.BalanceStr = b.String()

		if _token.PriceUSD != 0 {
			v := new(big.Float)
			v.Copy(b)
			v.Mul(v, big.NewFloat(_token.PriceUSD))

			_token.Value = *v
			_token.ValueStr = v.String()
		}
		result[chainId][_tokenId] = *_token

	}

}

func GetChainsTokenBalances(
	chainIds []schema.ChainId,
	wallet common.Address) map[schema.ChainId]schema.TokenMapping {

	chunkedResultChannel := make(chan []ChunkCall[*big.Int])
	_res := make(map[schema.ChainId]schema.TokenMapping)

	var totalChunkCount uint64
	totalChunkCount = 0

	for _, chainId := range chainIds {
		_tokens := configs.ChainTokens(chainId)
		_multicall := configs.ChainMultiCall(chainId)

		if _multicall == nil || _tokens == nil {
			continue
		}
		atomic.AddUint64(
			&totalChunkCount,
			getTokenBalances(TokenBalanceCallOpt, chainId, *_multicall, _tokens, wallet, chunkedResultChannel))
	}
	if totalChunkCount == 0 {
		return _res
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

	return _res
}

func GetChainsTokenAllowance(
	chainIds []schema.ChainId,
	wallet common.Address) map[schema.ChainId]schema.TokenMapping {

	chunkedResultChannel := make(chan []ChunkCall[*big.Int])
	_res := make(map[schema.ChainId]schema.TokenMapping)

	var totalChunkCount uint64
	totalChunkCount = 0

	for _, chainId := range chainIds {
		_tokens := configs.ChainTokens(chainId)
		_multicall := configs.ChainMultiCall(chainId)

		if _multicall == nil || _tokens == nil {
			continue
		}
		atomic.AddUint64(
			&totalChunkCount,
			getTokenBalances(TokenBalanceCallOpt, chainId, *_multicall, _tokens, wallet, chunkedResultChannel))
	}
	if totalChunkCount == 0 {
		return _res
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

	return _res
}
