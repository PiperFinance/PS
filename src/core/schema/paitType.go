package schema

import "math/big"

type PairId uint32

type PairDet struct {
	Tokens      map[TokenId]Token `json:"tokens"`
	TokensOrder []TokenId         `json:"tokensOrder"`
	Decimals    Decimals          `json:"decimals"`
	ChainId     ChainId           `json:"chainId"`
	Address     string            `json:"address"`
	Symbol      string            `json:"symbol"`
	Name        string            `json:"name"`
	Dex         string            `json:"dex"`
	Verified    bool              `json:"verified"`
	CoingeckoId string            `json:"coingeckoId,omitempty"`
}
type Pair struct {
	Detail PairDet `json:"detail"`
	//Reserves    []big.Int `json:"reserves"`
	//TotalSupply big.Int   `json:"totalSupply"`
	PriceUSD Price     `json:"priceUSD"`
	Balance  big.Int   `json:"balance"`
	Value    big.Float `json:"value"`
}

type PairMapping map[PairId]Pair
