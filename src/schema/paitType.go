package schema

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type PairId Id

type PairDet struct {
	Tokens      map[TokenId]Token `json:"tokens"`
	TokensOrder []TokenId         `json:"tokensOrder"`
	Decimals    Decimals          `json:"decimals"`
	ChainId     ChainId           `json:"chainId"`
	Address     common.Address    `json:"address"`
	Symbol      string            `json:"symbol"`
	Name        string            `json:"name"`
	Dex         string            `json:"dex"`
	Verified    bool              `json:"verified"`
	CoingeckoId string            `json:"coingeckoId,omitempty"`
}
type Pair struct {
	Detail PairDet `json:"detail"`
	// Reserves    []big.Int `json:"reserves"`
	// TotalSupply big.Int   `json:"totalSupply"`
	PriceUSD            float64   `json:"priceUSD"`
	Balance             big.Float `json:"-"`
	Value               big.Float `json:"-"`
	BalanceStr          string    `json:"balance"`
	BalanceNoDecimalStr string    `json:"balanceNoDecimal"`
	ValueStr            string    `json:"value"`
}

type PairMapping map[PairId]Pair

// Copy Returns another Object with same detail
func (pair Pair) Copy() *Pair {
	return &Pair{
		Detail: pair.Detail,
	}
}
