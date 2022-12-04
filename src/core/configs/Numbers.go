package configs

import (
	"math/big"
	"portfolio/core/schema"
)

func TEN() *big.Int {
	return big.NewInt(10)
}
func EIGHT() *big.Int {
	return big.NewInt(8)
}
func ZERO() *big.Int {
	return big.NewInt(0)
}
func ONE() *big.Int {
	return big.NewInt(0)
}

func init() {
	//_, _, _ = TEN, EIGHT, ZERO
}

func DecimalPowTen(decimals schema.Decimals) *big.Int {
	r := TEN()
	switch decimals {
	case 0:
		return ONE()
	case 8:
		r.Exp(r, EIGHT(), nil)
	default:
		r.Exp(r, big.NewInt(int64(decimals)), nil)
	}
	return r
}
