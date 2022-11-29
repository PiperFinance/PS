package configs

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	Multicall "portfolio/core/contracts/MulticallContract"

	"portfolio/core/schema"
	"sync"
)

var (
	onceForEthClient     sync.Once
	onceForMultiCall     sync.Once
	MULTICALL_V3_ADDRESS = common.HexToAddress("0xca11bde05977b3631167028862be2a173976ca11")
	gethClients          = make(map[schema.ChainId]*ethclient.Client, 100)
	multiCallInstances   = make(map[schema.ChainId]*Multicall.MulticallCaller, 100)
)

func EthClient(id schema.ChainId) *ethclient.Client {
	// TODO make it multi chain
	onceForEthClient.Do(func() {
		client, err := ethclient.Dial("https://cloudflare-eth.com")
		if err != nil {
			log.Fatal(err)
		}
		gethClients[1] = client
	})
	return gethClients[id]
}

func ChainMultiCall(id schema.ChainId) *Multicall.MulticallCaller {
	// TODO make it multi chain
	onceForMultiCall.Do(func() {
		contractInstance, err := Multicall.NewMulticallCaller(MULTICALL_V3_ADDRESS, EthClient(1))
		if err != nil {
			log.Fatal(err)
		}
		multiCallInstances[1] = contractInstance
	})
	return multiCallInstances[1]
}
