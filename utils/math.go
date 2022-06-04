package utils

import (
	"math"
	"math/big"

	"github.com/opefago/bot-tave/common"
	"github.com/shopspring/decimal"
)

func WeiToETH(balanceGwei *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(balanceGwei.String())
	return new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
}

/*
* Takes a starting token, an amount and a token-amount pair, and then calculates the amount either in the forward
* or reverse direction, outputing the receiving token*/
func CalculateTokenReturns(baseToken string, amount decimal.Decimal, priceInfo common.PriceInfo) (decimal.Decimal, string, common.Direction) {
	swap := decimal.NewFromInt(0)
	boughtToken := ""
	direction := common.Forward
	if baseToken == priceInfo.Pair.Base {
		swap = decimal.NewFromInt(1).Div(priceInfo.Price.Ask)
		boughtToken = priceInfo.Pair.Quote
		direction = common.Forward
	} else if baseToken == priceInfo.Pair.Quote {
		swap = priceInfo.Price.Bid
		boughtToken = priceInfo.Pair.Base
		direction = common.Reverse
	}
	return swap.Mul(amount), boughtToken, direction
}
