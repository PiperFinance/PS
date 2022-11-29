package configs

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"portfolio/core/schema"
	"sync"
)

var (
	onceTokenAddress sync.Once
	onceForMap       sync.Once
	once             sync.Once
	tokens           []schema.Token
	tokensAddress    []common.Address
	tokenChainMap    = make(map[schema.ChainId][]*schema.Token, 10)
)

func GetTokens() []schema.Token {
	once.Do(func() {
		// Load Tokens ...
		jsonFile, err := os.Open("core/data/tokens.json")
		defer jsonFile.Close()
		if err != nil {
			log.Fatalf("TokenLoader: %s", err)
		}
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatalf("TokenLoader: %s", err)
		}
		err = json.Unmarshal(byteValue, &tokens)
		if err != nil {
			log.Fatalf("TokenLoader: %s", err)
		}
	})
	return tokens
}

func GetTokensAddress() []common.Address {
	onceTokenAddress.Do(func() {
		// Load Tokens ...
		tokens := GetTokens()
		for _, t := range tokens {
			tokensAddress = append(tokensAddress, t.Get())
		}
	})
	return tokensAddress
}

func GetChainTokens(id schema.ChainId) []*schema.Token {
	onceForMap.Do(func() {
		// Load Tokens Chain Mapping
		jsonFileMap, err := os.Open("core/data/tokensChainMap.json")
		defer jsonFileMap.Close()
		if err != nil {
			log.Fatalf("TokenLoader: %s", err)
		}
		byteValueMap, _ := ioutil.ReadAll(jsonFileMap)
		err = json.Unmarshal(byteValueMap, &tokenChainMap)
		if err != nil {
			log.Fatalf("TokenLoader: %s", err)
		}
	})
	return tokenChainMap[id]
}
