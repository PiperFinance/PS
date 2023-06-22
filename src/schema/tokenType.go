package schema

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type TokenId Id

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
	Detail              TokenDet  `json:"detail"`
	PriceUSD            float64   `json:"priceUSD"`
	Balance             big.Float `json:"-"`
	Value               big.Float `json:"-"`
	BalanceStr          string    `json:"balance"`
	BalanceNoDecimalStr string    `json:"balanceNoDecimal"`
	ValueStr            string    `json:"value"`
}

// TODO make this a pointer !
type TokenMapping map[TokenId]Token

// Copy Only copies the detail bit
func (token Token) Copy() *Token {
	return &Token{
		Detail: token.Detail,
	}
}

func (token *Token) Id() TokenId {
	return TokenId(fmt.Sprintf("%s-%d", strings.ToLower(token.Detail.Address.String()), token.Detail.ChainId))
}

//func (tokenMapping TokenMapping) copy(src TokenMapping) TokenMapping {
//	dest := make(TokenMapping)
//	for k, v := range src {
//		dest[k] = v
//	}
//	return src
//}
