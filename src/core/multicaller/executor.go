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

// func execute(chunkIndex int, id schema.ChainId, wallet common.Address, multiCaller Multicall.MulticallCaller, chunkedCalls []ChunkCall[*big.Int], chunkChannel chan []ChunkCall[*big.Int]) {
// 	executeWithTimeout(
// 		chunkIndex, id, wallet, multiCaller, chunkedCalls,
// 		chunkChannel,
// 		configs.ChainContextTimeOut(id))
// }

func execute(chunkIndex int, id schema.ChainId, wallet common.Address, multiCaller Multicall.MulticallCaller, chunkedCalls []ChunkCall[*big.Int], chunkChannel chan []ChunkCall[*big.Int], callTimeout time.Duration) {
	calls := make([]Multicall.Multicall3Call3, len(chunkedCalls))
	for i, indexedCall := range chunkedCalls {
		calls[i] = indexedCall.Call
	}

	ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
	defer cancel()
	DefaultW3CallOpts := bind.CallOpts{Context: ctx}
	for i, _call := range calls {
		log.Infof("[%d][%s][%s]", i, _call.Target, common.Bytes2Hex(_call.CallData))
	}

	res, err := multiCaller.Aggregate3(&DefaultW3CallOpts, calls)

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
