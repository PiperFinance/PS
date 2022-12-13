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
	//onceTokenAddress sync.Once
	//tokensAddress    []common.Address

	onceForChainTokens sync.Once
	//tokens          []schema.Token
	// CD chain Tokens URL
	chainTokens          map[schema.ChainId]schema.TokenMapping
	NULL_TOKEN_ADDRESS   = common.HexToAddress("0x0000000000000000000000000000000000000000")
	NATIVE_TOKEN_ADDRESS = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	//tokenChainMap = make(map[schema.ChainId][]*schema.Token, 10)
)

func init() {

	onceForChainTokens.Do(func() {
		// Load Tokens ...
		// TODO READ FROM ENV
		tokensDir := "core/data/chain_separated_v2.json"
		chainTokensUrl := "https://github.com/PiperFinance/CD/blob/main/tokens/outVerified/chain_separated_v2.json?raw=true"
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
		err := json.Unmarshal(byteValue, &chainTokens)
		if err != nil {
			log.Fatalf("TokenLoader: %s", err)
		}
	})

}

func AllChainsTokens() []schema.ChainToken {
	return chainTokens
}

// TokensAddress Results in token object of this token address ...
func TokensAddress(id schema.ChainId, address common.Address) *schema.Token {
	for _, chainToken := range chainTokens {
		if chainToken.ChainId == id {
			for _, token := range chainToken.Tokens {
				if token.Address == address {
					return &token
				}
			}
		}
	}
	return nil
}

func ChainTokens(id schema.ChainId) []schema.Token {
	for _, chainToken := range chainTokens {
		if chainToken.ChainId == id {
			return chainToken.Tokens
		}
	}
	return nil
}
