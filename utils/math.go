package utils

import (
	"math"
	"math/big"
)

func WeiToETH(balanceGwei *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(balanceGwei.String())
	return new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
}
