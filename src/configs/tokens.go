package configs

import (
	"encoding/json"
	"errors"
	"fmt"
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
	allTokensArray       = make([]schema.Token, 0)
	allTokens            = make(schema.TokenMapping)
	chainTokens          = make(map[schema.ChainId]schema.TokenMapping)
	NULL_TOKEN_ADDRESS   = common.HexToAddress("0x0000000000000000000000000000000000000000")
	NATIVE_TOKEN_ADDRESS = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	tokensUrl            = "https://github.com/PiperFinance/CD/blob/main/tokens/outVerified/all_tokens.json?raw=true"
	tokensDir            = "data/all_tokens.json"
	priceUpdaterLock     = false
	priceUpdaterTTL      = 2 * time.Minute
	accessedChains       = ttlcache.New[string, []schema.ChainId](
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
			allTokensArray = append(allTokensArray, token)
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
	if priceUpdaterLock {
		return
	}
	priceUpdaterLock = true

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 5000
	t.MaxConnsPerHost = 1
	t.MaxIdleConnsPerHost = 5000

	httpClient := &http.Client{
		Timeout:   1 * time.Minute,
		Transport: t,
	}

	_chains := accessedChains.Get("ChainsToUpdate")
	if _chains == nil || _chains.IsExpired() {
		return
	}
	for _, chainId := range _chains.Value() {
		for tokenId, _ := range ChainTokens(chainId) {
			go func(id schema.TokenId, chainId schema.ChainId) {
				var tokenPrice float64
				res, err := httpClient.Get(fmt.Sprintf("%s?tokenId=%d", TpServer, id))
				if err != nil {
					log.Error(err)
				} else {
					byteValue, err := ioutil.ReadAll(res.Body)
					parseErr := json.Unmarshal(byteValue, &tokenPrice)
					if err != nil {
						log.Error(err)
					} else if parseErr != nil {
						log.Error(parseErr)
					} else {
						log.Infof("ID : %s  => %s", id, tokenPrice)
					}

				}
			}(tokenId, chainId)
		}
	}
	priceUpdaterLock = false

}

func AllChainsTokens() schema.TokenMapping {
	return allTokens
}
func AllChainsTokensArray() []schema.Token {
	return allTokensArray
}

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
