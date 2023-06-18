package configs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"portfolio/schema"
)

var (
	allPairs      = make(schema.PairMapping)
	chainPairs    = make(map[schema.ChainId]schema.PairMapping)
	allPairsArray = make([]schema.Pair, 0)
	pairsUrl      = "https://raw.githubusercontent.com/PiperFinance/CD/main/pair/all_pairs.json"
	pairsDir      = "data/all_pairs.json"
)

func LoadPairs() {
	// Load Pairs ...
	var byteValue []byte
	if _, err := os.Stat(pairsDir); errors.Is(err, os.ErrNotExist) {
		resp, err := http.Get(pairsUrl)
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
}

func AllChainsPairs() schema.PairMapping {
	return allPairs
}

func AllChainsPairsArray() schema.PairMapping {
	return allPairs
}

func ChainPairs(id schema.ChainId) schema.PairMapping {
	return chainPairs[id]
}

func ChainPairsArray(id schema.ChainId) schema.PairMapping {
	return chainPairs[id]
}
