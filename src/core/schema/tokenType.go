package schema

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type TokenId uint32

type TokenDet struct {
	ChainId     ChainId        `json:"chainId"`
	Address     common.Address `json:"address"`
	Name        string         `json:"name"`
	Symbol      string         `json:"symbol"`
	Decimals    Decimals       `json:"decimals"`
	Tags        []string       `json:"tags"`
	CoingeckoId string         `json:"coingeckoId"`
	LifiId      string         `json:"lifiId,omitempty"`
	ListedIn    []string       `json:"listedIn"`
	LogoURI     string         `json:"logoURI"`
	Verify      bool           `json:"verify"`
	Related     []Token        `json:"token,omitempty"`
}

type Token struct {
	Detail     TokenDet  `json:"detail"`
	PriceUSD   float64   `json:"priceUSD"`
	Balance    big.Float `json:"-"`
	Value      big.Float `json:"-"`
	BalanceStr string    `json:"balance"`
	ValueStr   string    `json:"value"`
}

type TokenMapping map[TokenId]Token

func (tokenMapping TokenMapping) copy(src TokenMapping) TokenMapping {
	dest := make(TokenMapping)
	for k, v := range src {
		dest[k] = v
	}
	return src
}
