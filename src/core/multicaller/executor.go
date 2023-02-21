package multicaller

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math/big"
	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/schema"
)

// what is either p / t > pair / token
func execute(what string, chunkIndex int, id schema.ChainId, wallet common.Address, multiCaller Multicall.MulticallCaller, chunkedCalls []ChunkCall[*big.Int], chunkChannel chan []ChunkCall[*big.Int]) {

	calls := make([]Multicall.Multicall3Call3, len(chunkedCalls))
	for i, indexedCall := range chunkedCalls {
		calls[i] = indexedCall.Call
	}

	contx, cancel := context.WithTimeout(context.Background(), configs.ChainContextTimeOut(id))
	defer cancel()
	DefaultW3CallOpts := bind.CallOpts{Context: contx}

	res, err := multiCaller.Aggregate3(&DefaultW3CallOpts, calls)
	//res2, err2 := multiCaller.GetChainId(&bind.CallOpts{Context: contx})
	//_, _ = res2, err2
	cacheKey := ChunkCallsCacheKey{wallet, id, chunkIndex, what}
	c := ChunkCallsCache.Get(cacheKey)
	if c != nil && !c.IsExpired() {
		ChunkCallsCache.Delete(cacheKey)
	}
	if err != nil {
		log.Error(err)
		for i, _ := range chunkedCalls {
			chunkedCalls[i].Err = err
		}
	} else {
		for i, _res := range res {
			chunkedCalls[i].CallRes = _res
			if _res.Success {
				_parser := chunkedCalls[i].ResultParser
				if _parser != nil {
					chunkedCalls[i].ParsedCallRes = chunkedCalls[i].ResultParser(_res.ReturnData)
				}
			} else {
				fmt.Println(_res.ReturnData)
			}
		}
		ChunkCallsCache.Set(cacheKey, chunkedCalls, configs.ChunkCallCacheTTL)
	}
	chunkChannel <- chunkedCalls
}
