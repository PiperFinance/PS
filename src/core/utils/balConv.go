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
	_decimal := configs.DecimalPowTen(token.Detail.Decimals)
	b := new(big.Float).SetInt(bal)
	token.BalanceNoDecimalStr = bal.String()
	b = b.Quo(b, new(big.Float).SetInt(_decimal))
	if b.IsInf() {
		return fmt.Errorf("inf bal")
	}
	token.Balance = *b
	token.BalanceStr = b.String()
	if token.PriceUSD != 0 {
		v := new(big.Float)
		v.Copy(b)
		v.Mul(v, big.NewFloat(token.PriceUSD))

		token.Value = *v
		token.ValueStr = v.String()
	}
	return nil
}
