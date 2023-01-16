package multicaller

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	log "github.com/sirupsen/logrus"
	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/schema"
)

func execute[T any](id schema.ChainId, multiCaller Multicall.MulticallCaller, chunkedCalls []ChunkCall[T], chunkChannel chan []ChunkCall[T]) {

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
	}
	chunkChannel <- chunkedCalls
}
