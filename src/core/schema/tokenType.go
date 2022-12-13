package schema

import "math/big"

type TokenId string

type TokenDet struct {
	ChainId     int      `json:"chainId"`
	Address     string   `json:"address"`
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
	Decimals    int      `json:"decimals"`
	Tags        []string `json:"tags"`
	CoingeckoId string   `json:"coingeckoId"`
	LifiId      string   `json:"lifiId,omitempty"`
	ListedIn    []string `json:"listedIn"`
	LogoURI     string   `json:"logoURI"`
	Verify      bool     `json:"verify"`
}

type Token struct {
	Token    TokenDet  `json:"token"`
	PriceUSD float64   `json:"priceUSD"`
	Balance  big.Int   `json:"balance"`
	Value    big.Float `json:"value"`
}

type TokenMapping map[TokenId]Token
