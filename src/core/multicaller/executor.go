package multicaller

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math/big"
	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/schema"
)

func execute(chunkIndex int, id schema.ChainId, wallet common.Address, multiCaller Multicall.MulticallCaller, chunkedCalls []ChunkCall[*big.Int], chunkChannel chan []ChunkCall[*big.Int]) {

	calls := make([]Multicall.Multicall3Call3, len(chunkedCalls))
	for i, indexedCall := range chunkedCalls {
		calls[i] = indexedCall.Call
	}

	contx, cancle := context.WithTimeout(context.Background(), configs.ChainContextTimeOut(id))
	defer cancle()
	DefaultW3CallOpts := bind.CallOpts{Context: contx}

	res, err := multiCaller.Aggregate3(&DefaultW3CallOpts, calls)

	cacheKey := ChunkCallsCacheKey{wallet, id, chunkIndex}
	ChunkCallsCache.Delete(cacheKey)
	if err != nil {
		log.Error(err)
		for i, _ := range chunkedCalls {
			chunkedCalls[i].Err = err
		}
	} else {
		for i, _res := range res {
			chunkedCalls[i].CallRes = _res
			if _res.Success {
				chunkedCalls[i].ParsedCallRes = chunkedCalls[i].ResultParser(_res.ReturnData)
			}
		}
		ChunkCallsCache.Set(cacheKey, chunkedCalls, ChunkCallCacheTTL)
	}
	chunkChannel <- chunkedCalls
}
