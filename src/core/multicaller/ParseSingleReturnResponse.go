package multicaller

import (
	"math/big"
	"portfolio/configs"
)

// ParseBigIntResult If there is only one bigint in call's response
func ParseBigIntResult(result []byte) *big.Int {
	z := configs.ZERO()
	if len(result) > 32 {
		z.SetBytes(result[:32])
	} else {
		z.SetBytes(result)
	}
	if z.Cmp(configs.MIN_BALANCE()) <= 0 {
		return nil
	} else {
		return z
	}
}
