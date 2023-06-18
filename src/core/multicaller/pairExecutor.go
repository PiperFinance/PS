package multicaller

import (
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/core/utils"
	"portfolio/schema"
)

var PairBalanceCallOpt ChunkedCallOpts

func init() {
	PairBalanceCallOpt = ChunkedCallOpts{W3CallOpt: nil, ChunkSize: 100, MaxRetries: 3, MaxTimeout: 1 * time.Minute}
}

// getPairBalances Wallet balance based on given pair ( Faster if chunks is used)
// Does not sort + only respond with pairs with balance
func getPairBalances(
	callOpts ChunkedCallOpts,
	id schema.ChainId,
	multiCaller *Multicall.MulticallCaller,
	pairs schema.PairMapping,
	wallet common.Address,
	chunkedResultChannel chan []ChunkCall[*big.Int],
) uint64 {
	allCalls := genPairBalanceCalls(pairs, wallet)
	chunkedCalls := utils.Chunks[ChunkCall[*big.Int]](allCalls, callOpts.ChunkSize)

	for i, indexCalls := range chunkedCalls {
		cachedChunkCalls := ChunkCallsCache.Get(ChunkCallsCacheKey{wallet, id, i})
		if cachedChunkCalls != nil && !cachedChunkCalls.IsExpired() {
			go func() {
				chunkedResultChannel <- cachedChunkCalls.Value()
			}()
		} else {
			go executeWithRetries(i, id, wallet, multiCaller, indexCalls, chunkedResultChannel, callOpts.MaxTimeout, callOpts.MaxRetries)
		}
	}

	return uint64(len(chunkedCalls))
}

func balancePairResultParser(wallet common.Address, result map[schema.ChainId]schema.PairMapping, chunk []ChunkCall[*big.Int]) {
	for _, call := range chunk {

		if !call.CallRes.Success || call.ParsedCallRes == nil {
			continue
		}
		chainId := call.ChainId
		if result[chainId] == nil {
			result[chainId] = make(schema.PairMapping)
		}
		pairId := schema.PairId(call.Id)

		_pair := configs.ChainPairs(chainId)[pairId].Copy()

		// Pair Balance
		decimal := configs.DecimalPowTen(_pair.Detail.Decimals)
		b := new(big.Float).SetInt(call.ParsedCallRes)
		b = b.Quo(b, new(big.Float).SetInt(decimal))
		if b.IsInf() {
			log.Errorf("[INF] @ (%d,%s) ::: cnResp  %s ", _pair.Detail.Decimals, _pair.Detail.Address, call.ParsedCallRes)
		}
		_pair.Balance = *b
		_pair.BalanceStr = b.String()
		_pair.BalanceNoDecimalStr = call.ParsedCallRes.String()

		if _pair.PriceUSD != 0 {
			v := new(big.Float)
			v.Copy(b)
			v.Mul(v, big.NewFloat(_pair.PriceUSD))

			_pair.Value = *v
			_pair.ValueStr = v.String()
		}
		result[chainId][pairId] = *_pair

	}
	_ = wallet
}

func GetChainsPairBalances(
	chainIds []schema.ChainId,
	wallet common.Address,
) (map[schema.ChainId]schema.PairMapping, error) {
	chunkedResultChannel := make(chan []ChunkCall[*big.Int])
	_res := make(map[schema.ChainId]schema.PairMapping)

	var totalChunkCount uint64
	// totalChunkCount = 0

	for _, chainId := range chainIds {
		_pairs := configs.ChainPairs(chainId)
		_multicall, err := configs.ChainMultiCall(chainId)
		if err != nil {
			return nil, err
		}

		if _multicall == nil || _pairs == nil {
			continue
		}
		atomic.AddUint64(
			&totalChunkCount,
			getPairBalances(PairBalanceCallOpt, chainId, _multicall, _pairs, wallet, chunkedResultChannel))
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

	return _res, nil
}
