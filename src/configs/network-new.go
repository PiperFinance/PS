package configs

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/schema"
	"portfolio/utils"
)

var (
	// EthClient  *ethclient.Client
	EthClientS           map[int64][]*ethclient.Client
	selectorMutex        sync.Mutex
	selectorIndex        map[int64]int
	clientCount          map[int64]int
	rpcs                 map[int64][]string
	MULTICALL_V3_ADDRESS = common.HexToAddress("0xca11bde05977b3631167028862be2a173976ca11")
	Networks             = make([]schema.Network, 0)
)

func LoadNetwork() {
	selectorMutex = sync.Mutex{}
	rpcs = make(map[int64][]string, len(Config.SupportedChains))
	EthClientS = make(map[int64][]*ethclient.Client, len(Config.SupportedChains))
	clientCount = make(map[int64]int, len(Config.SupportedChains))
	selectorIndex = make(map[int64]int, len(Config.SupportedChains))
	for _, net := range SupportedNetworks {
		rpcs[net.ChainId] = utils.GetNetworkRpcUrls(net.GoodRpc)
		clientCount[net.ChainId] = len(net.GoodRpc)
		selectorIndex[net.ChainId] = 0
		EthClientS[net.ChainId] = make([]*ethclient.Client, len(net.GoodRpc))
		for i, _rpc := range net.GoodRpc {

			client, err := ethclient.Dial(_rpc.Url)
			if err != nil {
				Logger.Panicf("Client Connection %+v Error : %s  ", _rpc, err)
			}
			EthClientS[net.ChainId][i] = client
		}
	}
}

func ChainContextTimeOut(id schema.ChainId) time.Duration {
	return time.Millisecond * 5450
}

func ChainMultiCall(id schema.ChainId) (*Multicall.MulticallCaller, error) {
	contractInstance, err := Multicall.NewMulticallCaller(MULTICALL_V3_ADDRESS, EthClient(int64(id)))
	if err != nil {
		// log.Errorf("Contract Instance Creation Error : %s @ chainID :%d", err, id)
		return nil, err
	}
	return contractInstance, nil
	// multiCallInstances[chainId] = contractInstance
	// return multiCallInstances[id]
}

func EthClient(chain int64) *ethclient.Client {
	cl, _ := EthClientDebug(chain)
	return cl
}

func EthClientDebug(chain int64) (*ethclient.Client, string) {
	defer func() {
		selectorMutex.Lock()
		selectorIndex[chain]++
		if selectorIndex[chain] >= clientCount[chain] {
			selectorIndex[chain] = 0
		}
		selectorMutex.Unlock()
	}()
	// TODO - Try to recover panic here !
	if clients, ok := EthClientS[chain]; ok {
		if index, ok := selectorIndex[chain]; ok {
			if len(clients) == 0 {
				Logger.Errorf("EthCLient Selector : No RPCs found for chain %d", chain)
			} else {
				return clients[index], rpcs[chain][index]
			}
		}
	}
	return nil, ""
}

// BatchLogMaxHeight staticky returns block height set in mainnet.json !
func BatchLogMaxHeight(chain int64) uint64 {
	r := uint64(SupportedNetworks[chain].BatchLogMaxHeight)
	if r == 0 {
		// TODO make this dynamic
		return 1
	} else {
		return r
	}
}

// BatchLogMaxHeight staticky returns block height set in mainnet.json !
func MulticallMaxSize(chain int64) uint64 {
	r := uint64(SupportedNetworks[chain].MulticallMaxSize)
	if r == 0 {
		// TODO make this dynamic
		return 1
	} else {
		return r
	}
}

func MulticallAddress(chain int64) common.Address {
	return MULTICALL_V3_ADDRESS
}

// LatestBlock Last block mines head delay for safe data aggregation (uncle blocks!)
func LatestBlock(ctx context.Context, chain int64) (uint64, error) {
	if b, err := EthClient(chain).BlockNumber(ctx); err != nil {
		return Config.StartingBlockNumber, err
	} else {
		return b - Config.BlockHeadDelay, nil
	}
}

func NetworkValueAddress(chain int64) common.Address {
	return MULTICALL_V3_ADDRESS
}
