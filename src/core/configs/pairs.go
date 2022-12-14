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
	allPairs          schema.TokenMapping
	chainPairs        map[schema.ChainId]schema.TokenMapping
	allPairsArray     []schema.Token
	chainPairsUrl     = "https://raw.githubusercontent.com/PiperFinance/CD/main/pair/all_pairs.json"
	pairsDir          = "core/data/all_tokens.json"
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
			chainTokens[schema.ChainId(token.Token.ChainId)][tokenId] = token
			allTokensArray = append(allTokensArray, token)
		}

	})

}

func AllChainsTokens() schema.TokenMapping {
	return allTokens
}

func ChainTokens(id schema.ChainId) schema.TokenMapping {
	return chainTokens[id]
}
