package configs

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"portfolio/core/schema"
	"sync"
)

var (
	onceForChainPairs sync.Once
	allPairs          = make(schema.PairMapping)
	chainPairs        = make(map[schema.ChainId]schema.PairMapping)
	allPairsArray     = make([]schema.Pair, 0)
	chainPairsUrl     = "https://raw.githubusercontent.com/PiperFinance/CD/main/pair/all_pairs.json"
	pairsDir          = "core/data/all_pairs.json"
)

func init() {

	onceForChainPairs.Do(func() {
		// Load Tokens ...
		// TODO READ FROM ENV
		var byteValue []byte
		if _, err := os.Stat(pairsDir); errors.Is(err, os.ErrNotExist) {
			resp, err := http.Get(chainPairsUrl)
			if err != nil {
				log.Fatalln(err)
			}
			byteValue, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("HTTPPairLoader: %s", err)
			}
		} else {
			jsonFile, err := os.Open(pairsDir)
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
		err := json.Unmarshal(byteValue, &allPairs)
		if err != nil {
			log.Fatalf("PairLoader: %s", err)
		}

		for pairId, pair := range allPairs {
			chainId := pair.Detail.ChainId
			if chainPairs[chainId] == nil {
				chainPairs[chainId] = make(schema.PairMapping)
			}
			chainPairs[chainId][pairId] = pair
			allPairsArray = append(allPairsArray, pair)
		}

	})

}

func AllChainsPairs() schema.TokenMapping {
	return allTokens
}

func AllChainsPairsArray() schema.TokenMapping {
	return allTokens
}

func ChainPairs(id schema.ChainId) schema.TokenMapping {
	return chainTokens[id]
}
func ChainPairsArray(id schema.ChainId) schema.TokenMapping {
	return chainTokens[id]
}
