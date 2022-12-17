package multicaller

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	BalanceCallOpt ChunkedCallOpts
)

func init() {
	BalanceCallOpt = ChunkedCallOpts{W3CallOpt: nil, ChunkSize: 250}
}

// GetBalancesFaster Wallet balance based on given token ( Faster if chunks is used)
// Does not sort + only respond with tokens with balance
func GetBalancesFaster(
	callOpts ChunkedCallOpts,
	id schema.ChainId,
	multiCaller Multicall.MulticallCaller,
	tokens schema.TokenMapping,
	wallets common.Address,
	chunkedResultChannel chan []chunkCall[*big.Int]) uint64 {

	allCalls := genGetBalanceCalls(tokens, wallets)
	chunkedCalls := utils.Chunks[chunkCall[*big.Int]](allCalls, callOpts.ChunkSize)

	for _, indexCalls := range chunkedCalls {

		go func(chunkedCalls []chunkCall[*big.Int], chunkChannel chan []chunkCall[*big.Int]) {

			calls := make([]Multicall.Multicall3Call3, len(chunkedCalls))
			for i, indexedCall := range chunkedCalls {
				calls[i] = indexedCall.Call
			}

			contx, cancle := context.WithTimeout(context.Background(), configs.ChainContextTimeOut(id))
			defer cancle()
			DefaultW3CallOpts := bind.CallOpts{Context: contx}

			res, err := multiCaller.Aggregate3(&DefaultW3CallOpts, calls)

			if err != nil {
				log.Error(err)
				chunkChannel <- nil
			} else {
				for i, _res := range res {
					chunkedCalls[i].CallRes = _res
					if _res.Success {
						chunkedCalls[i].ParsedCallRes = chunkedCalls[i].ResultParser(_res.ReturnData)
					} else {
						chunkedCalls[i].ParsedCallRes = nil
					}
				}
				chunkChannel <- chunkedCalls
			}
		}(indexCalls, chunkedResultChannel)
	}

	return uint64(len(chunkedCalls))
}

func GetChainsBalances(
	chainIds []schema.ChainId,
	wallet common.Address) map[schema.ChainId]schema.TokenMapping {
	chunkedResultChannel := make(chan []chunkCall[*big.Int])
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
			GetBalancesFaster(BalanceCallOpt, chainId, *_multicall, _tokens, wallet, chunkedResultChannel))
	}

	for tmp := range chunkedResultChannel {
		//tmp := <-chunkedResultChannel
		if totalChunkCount > 0 {
			totalChunkCount--
		}

		for _, call := range tmp {

			if !call.CallRes.Success || call.ParsedCallRes == nil {
				continue
			}
			chainId := call.ChainId
			if _res[chainId] == nil {
				_res[chainId] = make(schema.TokenMapping)
			}

			_token := configs.ChainTokens(chainId)[call.TokenId].Copy()

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
			_res[chainId][call.TokenId] = *_token

		}

		if totalChunkCount == 0 {
			break
		}
	}

	close(chunkedResultChannel)

	return _res
}
