package utils

import (
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"

	"portfolio/configs"
	"portfolio/schema"
)

// TODO - make this generic !
func MustParseBal(bal *big.Int, token *schema.Token) {
	if err := ParseBal(bal, token); err != nil {
		logrus.Errorf("[INF] @ (%d,%s) ::: cnResp  %s ", token.Detail.Decimals, token.Detail.Address, bal)
	}
}

func ParseBalAndParse(bal string, token *schema.Token) error {
	v := big.Int{}
	b, ok := v.SetString(bal, 10)
	if ok {
		return ParseBal(b, token)
	} else {
		return fmt.Errorf("bal %s to big.Int conv failed", bal)
	}
}

func ParseBal(bal *big.Int, token *schema.Token) error {
	prevBal := new(big.Float)
	prevBal.SetString(token.BalanceNoDecimalStr)

	_decimal := configs.DecimalPowTen(token.Detail.Decimals)

	b := new(big.Float).SetInt(bal)

	b.Add(b, prevBal)
	i, _ := b.Int(big.NewInt(0))

	token.BalanceNoDecimalStr = i.String()

	b = b.Quo(b, new(big.Float).SetInt(_decimal))
	if b.IsInf() {
		return fmt.Errorf("inf bal")
	}
	token.Balance = *b
	token.BalanceStr = b.String()
	z := configs.AllTokens[token.Id()]
	if z.PriceUSD != 0 {
		v := new(big.Float)
		v.Copy(b)
		v.Mul(v, big.NewFloat(z.PriceUSD))
		token.PriceUSD = z.PriceUSD
		token.Value = *v
		token.ValueStr = v.String()
	}
	return nil
}
