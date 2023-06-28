package utils

import (
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"

	"portfolio/configs"
	"portfolio/schema"
)

// TODO - make this generic !
func MustParseBalPair(bal *big.Int, token *schema.Pair) {
	if err := ParseBalPair(bal, token); err != nil {
		logrus.Errorf("[INF] @ (%d,%s) ::: cnResp  %s ", token.Detail.Decimals, token.Detail.Address, bal)
	}
}

func ParseBalAndParsePair(bal string, token *schema.Pair) error {
	v := big.Int{}
	b, ok := v.SetString(bal, 10)
	if ok {
		return ParseBalPair(b, token)
	} else {
		return fmt.Errorf("bal %s to big.Int conv failed", bal)
	}
}

func ParseBalPair(bal *big.Int, pair *schema.Pair) error {
	prevBal := new(big.Float)
	prevBal.SetString(pair.BalanceNoDecimalStr)

	_decimal := configs.DecimalPowTen(pair.Detail.Decimals)

	b := new(big.Float).SetInt(bal)

	b.Add(b, prevBal)
	i, _ := b.Int(big.NewInt(0))

	pair.BalanceNoDecimalStr = i.String()

	b = b.Quo(b, new(big.Float).SetInt(_decimal))
	if b.IsInf() {
		return fmt.Errorf("inf bal")
	}
	pair.Balance = *b
	pair.BalanceStr = b.String()
	z := configs.AllPairs[pair.Id()]
	if z.PriceUSD != 0 {
		v := new(big.Float)
		v.Copy(b)
		v.Mul(v, big.NewFloat(z.PriceUSD))
		pair.PriceUSD = z.PriceUSD
		pair.Value = *v
		pair.ValueStr = v.String()
	}
	return nil
}
