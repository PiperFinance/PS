package configs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/schema"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

var (
	onceForMainNet       sync.Once
	Networks             = make([]schema.Network, 0)
	MULTICALL_V3_ADDRESS = common.HexToAddress("0xca11bde05977b3631167028862be2a173976ca11")
	GethClients          = make(map[schema.ChainId]*ethclient.Client, 10)
	multiCallInstances   = make(map[schema.ChainId]*Multicall.MulticallCaller, 10)
	ChainIds             = make([]schema.ChainId, 0)
	chainsUrl            = "https://github.com/PiperFinance/CD/blob/main/chains/mainnet.json?raw=true"
	chainsDir            = "data/mainnet.json"
)

func init() {
	onceForMainNet.Do(func() {
		// Load Tokens ...

		var byteValue []byte
		if _, err := os.Stat(chainsDir); errors.Is(err, os.ErrNotExist) {
			resp, err := http.Get(chainsUrl)
			if err != nil {
				log.Fatalln(err)
			}
			byteValue, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("HTTPPairLoader: %s", err)
			}
		} else {
			jsonFile, err := os.Open(chainsDir)
			defer func(jsonFile *os.File) {
				err := jsonFile.Close()
				if err != nil {
					log.Error(err)
				}
			}(jsonFile)
			if err != nil {
				log.Fatalf("JSONPairLoader: %s", err)
			}
			byteValue, err = ioutil.ReadAll(jsonFile)
			if err != nil {
				log.Fatalf("JSONPairLoader: %s", err)
			}
		}
		err := json.Unmarshal(byteValue, &Networks)
		if err != nil {
			log.Fatalf("ChainsLoader: %s", err)
		}
		for _, chain := range Networks {
			chainId := schema.ChainId(chain.ChainId)
			client, err := ethclient.Dial(chain.RpcUrls.Default)
			if err != nil {
				log.Errorf("Client Connection Error : %s  @ chainId: %d", err, chainId)
			} else {
				GethClients[chainId] = client
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

func ChainContextTimeOut(id schema.ChainId) time.Duration {
	return time.Millisecond * 5450
}

func ChainMultiCall(id schema.ChainId) *Multicall.MulticallCaller {
	return multiCallInstances[id]
}
