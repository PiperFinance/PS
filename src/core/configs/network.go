package configs

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	Multicall "portfolio/core/contracts/MulticallContract"

	"portfolio/core/schema"
	"sync"
)

var (
	onceForEthClient     sync.Once
	onceForMultiCall     sync.Once
	onceForMainNet       sync.Once
	Networks             []schema.Network
	MULTICALL_V3_ADDRESS = common.HexToAddress("0xca11bde05977b3631167028862be2a173976ca11")
	gethClients          = make(map[schema.ChainId]*ethclient.Client, 10)
	multiCallInstances   = make(map[schema.ChainId]*Multicall.MulticallCaller, 10)
	ChainIds             []schema.ChainId
)

func init() {
	onceForMainNet.Do(func() {

		// Load Tokens ...
		jsonFile, err := os.Open("core/data/mainnet.json")
		defer func(jsonFile *os.File) {
			err := jsonFile.Close()
			if err != nil {
				log.Error(err)
			}
		}(jsonFile)
		if err != nil {
			log.Fatalf("ChainsLoader: %s", err)
		}
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatalf("ChainsLoader: %s", err)
		}
		err = json.Unmarshal(byteValue, &Networks)
		if err != nil {
			log.Fatalf("ChainsLoader: %s", err)
		}
		ChainIds = make([]schema.ChainId, 0)
		for _, chain := range Networks {
			chainId := schema.ChainId(chain.ChainId)
			client, err := ethclient.Dial(chain.RpcUrls.Default)
			if err != nil {
				log.Errorf("Client Connection Error : %s  @ chainId: %d", err, chainId)
			} else {
				gethClients[chainId] = client
				contractInstance, err := Multicall.NewMulticallCaller(MULTICALL_V3_ADDRESS, client)
				if err != nil {
					log.Errorf("Contract Instance Creation Error : %s @ chainID :%d", err, chainId)
				}
				multiCallInstances[chainId] = contractInstance
				ChainIds = append(ChainIds, chainId)
			}
		}
	})
}

func EthClient(id schema.ChainId) *ethclient.Client {
	// TODO make it multi chain
	//onceForEthClient.Do(func() {
	//	client, err := ethclient.Dial("https://cloudflare-eth.com")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	gethClients[1] = client
	//})
	return gethClients[id]
}

func ChainMultiCall(id schema.ChainId) *Multicall.MulticallCaller {
	// TODO make it multi chain
	//onceForMultiCall.Do(func() {
	//	contractInstance, err := Multicall.NewMulticallCaller(MULTICALL_V3_ADDRESS, EthClient(1))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	multiCallInstances[1] = contractInstance
	//})
	return multiCallInstances[id]
}
