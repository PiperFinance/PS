package configs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jellydator/ttlcache/v3"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"portfolio/schema"
)

var (
	// CD chain Tokens URL
	allTokensArray       = make([]schema.Token, 0)
	AllTokens            = make(schema.TokenMapping)
	ValueTokenIds        = make(map[schema.ChainId]schema.TokenId)
	ValueTokens          = make(map[schema.ChainId]schema.Token)
	chainTokens          = make(map[schema.ChainId]schema.TokenMapping)
	NULL_TOKEN_ADDRESS   = common.HexToAddress("0x0000000000000000000000000000000000000000")
	NATIVE_TOKEN_ADDRESS = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	tokensUrl            = "https://github.com/PiperFinance/CD/blob/main/tokens/outVerified/all_tokens.json?raw=true"
	tokensDir            = "data/all_tokens.json"
	priceUpdaterLock     = false
	priceUpdaterTTL      = 1 * time.Minute
	accessedChains       = ttlcache.New[string, []schema.ChainId](
		ttlcache.WithTTL[string, []schema.ChainId](15 * time.Second),
	)
)

func LoadTokens() {
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
	err := json.Unmarshal(byteValue, &AllTokens)
	if err != nil {
		log.Fatalf("TokenLoader: %s", err)
	}
	for tokenId, token := range AllTokens {
		chainId := token.Detail.ChainId
		if chainTokens[chainId] == nil {
			chainTokens[chainId] = make(schema.TokenMapping)
		}
		chainTokens[chainId][tokenId] = token
		allTokensArray = append(allTokensArray, token)
		if token.Detail.Address == common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE") {
			ValueTokenIds[chainId] = tokenId
			ValueTokens[chainId] = token
		}
	}

	cr := cron.New()
	priceUpdaterJobId, err := cr.AddFunc("*/1 * * * *", priceUpdater)
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("Started priceUpdaterJobId [%d] @ %+v", priceUpdaterJobId, time.Now())
	}
	cr.Start()
}

func priceUpdater() {
	// TODO - should be a mutex
	if priceUpdaterLock {
		return
	}
	priceUpdaterLock = true

	t := http.DefaultTransport.(*http.Transport).Clone()
	httpClient := &http.Client{
		Timeout:   1 * time.Minute,
		Transport: t,
	}
	resp := make(map[schema.TokenId]float64)
	res, _ := httpClient.Get(Config.TokenPriceURL.JoinPath("/all").String())
	if res != nil {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(body, &resp); err != nil {
			log.Error(err)
		} else {
			for tokenId, price := range resp {
				t, ok := AllTokens[tokenId]
				if ok {
					// AllTokens[tokenId].PriceUSD = price
					t.PriceUSD = price
					AllTokens[tokenId] = t
				}
			}
		}
	}

	priceUpdaterLock = false
}

func AllChainsTokens() schema.TokenMapping {
	return AllTokens
}

func AllChainsTokensArray() []schema.Token {
	return allTokensArray
}

func ChainTokens(id schema.ChainId) schema.TokenMapping {
	return chainTokens[id]
	// _chains := accessedChains.Get("ChainsToUpdate")
	// chains := make([]schema.ChainId, 1)
	// _ = chains
	// if _chains == nil {
	// 	chains = make([]schema.ChainId, 1)
	// 	chains[0] = id
	// } else {
	// 	chains = append(_chains.Value(), id)
	// }
	// accessedChains.Set("ChainsToUpdate", chains, priceUpdaterTTL)
	// t := chainTokens[id]
	// return t
}
