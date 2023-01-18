package configs

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jellydator/ttlcache/v3"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"net/http"
	"os"
	"portfolio/schema"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

var (
	TpServer           string
	onceForChainTokens sync.Once
	// CD chain Tokens URL
	//allTokensArray       = make([]schema.Token, 0)
	allTokens            = make(schema.TokenMapping)
	chainTokens          = make(map[schema.ChainId]schema.TokenMapping)
	NULL_TOKEN_ADDRESS   = common.HexToAddress("0x0000000000000000000000000000000000000000")
	NATIVE_TOKEN_ADDRESS = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	tokensUrl            = "https://github.com/PiperFinance/CD/blob/main/tokens/outVerified/all_tokens.json?raw=true"
	tokensDir            = "data/all_tokens.json"
	//priceUpdaterLock     = false
	priceUpdaterTTL = 15 * time.Minute
	accessedChains  = ttlcache.New[string, []schema.ChainId](
		ttlcache.WithTTL[string, []schema.ChainId](15 * time.Second),
	)
)

func init() {

	onceForChainTokens.Do(func() {

		// Load Tokens ...
		var byteValue []byte
		if _, err := os.Stat(tokensDir); errors.Is(err, os.ErrNotExist) {
			resp, err := http.Get(tokensUrl)
			if err != nil {
				log.Fatalln(err)
			}
			byteValue, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("HTTPTokenLoader: %s", err)
			}
		} else {
			jsonFile, err := os.Open(tokensDir)
			defer func(jsonFile *os.File) {
				err := jsonFile.Close()
				if err != nil {
					log.Error(err)
				}
			}(jsonFile)
			if err != nil {
				log.Fatalf("JSONTokenLoader: %s", err)
			}
			byteValue, err = ioutil.ReadAll(jsonFile)
			if err != nil {
				log.Fatalf("JSONTokenLoader: %s", err)
			}
		}
		err := json.Unmarshal(byteValue, &allTokens)
		if err != nil {
			log.Fatalf("TokenLoader: %s", err)
		}
		for tokenId, token := range allTokens {
			chainId := token.Detail.ChainId
			if chainTokens[chainId] == nil {
				chainTokens[chainId] = make(schema.TokenMapping)
			}
			chainTokens[chainId][tokenId] = token
			//allTokensArray = append(allTokensArray, token)
		}
	})
	cr := cron.New()
	priceUpdaterJobId, err := cr.AddFunc("*/2 * * * *", priceUpdater)
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("Started priceUpdaterJobId [%s] @ %s", priceUpdaterJobId, time.Now())
	}
	cr.Start()

	_TpServer, ok := os.LookupEnv("TP_SERVER")
	if !ok {
		_TpServer = "http://localhost:3001"
	}
	TpServer = _TpServer

}

func priceUpdater() {
	//if priceUpdaterLock {
	//	return
	//}
	//priceUpdaterLock = true

	_chains := accessedChains.Get("ChainsToUpdate")

	if _chains == nil || _chains.IsExpired() {
		return
	}
	ids := AllChainsTokenIds()
	res := make(map[schema.TokenId]float64)

	bytesValue, err := json.Marshal(ids)
	if err != nil {
		log.Error(err)
	}
	_res, err := http.Post(TpServer, "application/json", bytes.NewBuffer(bytesValue))
	if err != nil {
		log.Error(err)
		return
	}
	body, err := ioutil.ReadAll(_res.Body)
	if err != nil {
		log.Error(err)
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Error(err)
		return
	}
	for tokenId, price := range res {
		x := allTokens[tokenId]
		z := chainTokens[x.Detail.ChainId][tokenId]
		x.PriceUSD = price
		z.PriceUSD = price
		chainTokens[z.Detail.ChainId][tokenId] = z
		allTokens[tokenId] = x
		//log.Infof("ID : %s  => %s", tokenId, price)
	}

}

func AllChainsTokens() schema.TokenMapping {
	return allTokens
}
func AllChainsTokenIds() []schema.TokenId {
	allTokenIds := make([]schema.TokenId, 0)
	for tokenId, _ := range AllChainsTokens() {
		allTokenIds = append(allTokenIds, tokenId)
	}
	return allTokenIds
}

//func AllChainsTokensArray() []schema.Token {
//	return allTokensArray
//}

func ChainTokens(id schema.ChainId) schema.TokenMapping {
	_chains := accessedChains.Get("ChainsToUpdate")
	chains := make([]schema.ChainId, 1)
	if _chains == nil {
		chains = make([]schema.ChainId, 1)
		chains[0] = id
	} else {
		chains = append(_chains.Value(), id)
	}
	accessedChains.Set("ChainsToUpdate", chains, priceUpdaterTTL)
	t := chainTokens[id]
	return t
}
