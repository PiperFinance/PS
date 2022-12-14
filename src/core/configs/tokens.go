package configs

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"portfolio/core/schema"
	"sync"
)

var (
	onceForChainTokens sync.Once
	// CD chain Tokens URL
	allTokensArray       []schema.Token
	allTokens            schema.TokenMapping
	chainTokens          map[schema.ChainId]schema.TokenMapping
	NULL_TOKEN_ADDRESS   = common.HexToAddress("0x0000000000000000000000000000000000000000")
	NATIVE_TOKEN_ADDRESS = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	tokensDir            = "core/data/all_tokens.json"
	chainTokensUrl       = "https://github.com/PiperFinance/CD/blob/main/tokens/outVerified/all_tokens.json?raw=true"
)

func init() {

	onceForChainTokens.Do(func() {
		// Load Tokens ...
		// TODO READ FROM ENV
		var byteValue []byte
		if _, err := os.Stat(tokensDir); errors.Is(err, os.ErrNotExist) {
			resp, err := http.Get(chainTokensUrl)
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
