package multicaller

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/schema"
)

func executeWithRetries(chunkIndex int, id schema.ChainId, wallet common.Address, multiCaller *Multicall.MulticallCaller, chunkedCalls []ChunkCall[*big.Int], chunkChannel chan []ChunkCall[*big.Int], callTimeout time.Duration, retires int) {
	calls := make([]Multicall.Multicall3Call3, len(chunkedCalls))
	for i, indexedCall := range chunkedCalls {
		calls[i] = indexedCall.Call
	}

	var res []Multicall.Multicall3Result
	var err error
	for retires >= 0 {
		ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
		defer cancel()
		DefaultW3CallOpts := bind.CallOpts{Context: ctx}
		if _res, _err := multiCaller.Aggregate3(&DefaultW3CallOpts, calls); _err == nil {
			res = _res
			err = _err
			break
		} else {
			err = _err
			retires--
			continue
		}
	}
	cacheKey := ChunkCallsCacheKey{wallet, id, chunkIndex}
	ChunkCallsCache.Delete(cacheKey)
	if err != nil {
		log.Error(err)
		for i := range chunkedCalls {
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
