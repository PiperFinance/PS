package multicaller

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math/big"
	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/core/utils"
	"portfolio/schema"
	"sync/atomic"
)

var (
	PairBalanceCallOpt ChunkedCallOpts
)

func init() {
	PairBalanceCallOpt = ChunkedCallOpts{W3CallOpt: nil, ChunkSize: 250}
}

// getPairBalances Wallet balance based on given pair ( Faster if chunks is used)
// Does not sort + only respond with pairs with balance
func getPairBalances(
	callOpts ChunkedCallOpts,
	id schema.ChainId,
	multiCaller Multicall.MulticallCaller,
	pairs schema.PairMapping,
	wallets common.Address,
	chunkedResultChannel chan []ChunkCall[*big.Int]) uint64 {

	allCalls := genPairBalanceCalls(pairs, wallets)
	chunkedCalls := utils.Chunks[ChunkCall[*big.Int]](allCalls, callOpts.ChunkSize)

	for _, indexCalls := range chunkedCalls {
		go execute[*big.Int](id, multiCaller, indexCalls, chunkedResultChannel)
	}

	return uint64(len(chunkedCalls))
}

func balancePairResultParser(result map[schema.ChainId]schema.PairMapping, chunk []ChunkCall[*big.Int]) {
	for _, call := range chunk {
		if call.Err != nil {
			// TODO
			fmt.Println(call.Err)
		}
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

		if _pair.PriceUSD != 0 {
			v := new(big.Float)
			v.Copy(b)
			v.Mul(v, big.NewFloat(_pair.PriceUSD))

			_pair.Value = *v
			_pair.ValueStr = v.String()
		}
		result[chainId][pairId] = *_pair

	}

}

func GetChainsPairBalances(
	chainIds []schema.ChainId,
	wallet common.Address) map[schema.ChainId]schema.PairMapping {

	chunkedResultChannel := make(chan []ChunkCall[*big.Int])
	_res := make(map[schema.ChainId]schema.PairMapping)

	var totalChunkCount uint64
	totalChunkCount = 0

	for _, chainId := range chainIds {
		_pairs := configs.ChainPairs(chainId)
		_multicall := configs.ChainMultiCall(chainId)

		if _multicall == nil || _pairs == nil {
			continue
		}
		atomic.AddUint64(
			&totalChunkCount,
			getPairBalances(PairBalanceCallOpt, chainId, *_multicall, _pairs, wallet, chunkedResultChannel))
	}

	for chunkCalls := range chunkedResultChannel {
		//tmp := <-chunkedResultChannel
		if totalChunkCount > 0 {
			totalChunkCount--
		}

		balancePairResultParser(_res, chunkCalls)

		if totalChunkCount == 0 {
			break
		}
	}

	close(chunkedResultChannel)

	return _res
}
