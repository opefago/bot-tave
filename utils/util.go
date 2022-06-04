package utils

import (
	"github.com/opefago/bot-tave/common"
	"github.com/shopspring/decimal"
)

func FilterByStableCoins(traingularPair []common.TraingularPair, stableCoins []string) []common.TraingularPair {
	var sortedTraingularPair []common.TraingularPair
	for _, pair := range traingularPair {
		for _, coin := range stableCoins {
			if (coin == pair.PrimaryTickerPair.Base || coin == pair.PrimaryTickerPair.Quote) &&
				(coin == pair.TertiaryTickerPair.Base || coin == pair.TertiaryTickerPair.Quote) {
				sortedTraingularPair = append(sortedTraingularPair, pair)
			}
		}
	}
	return sortedTraingularPair
}

func CalculateRealPotential(startAmount decimal.Decimal, normalisedPrices []common.PriceQuantity) decimal.Decimal {
	balance := startAmount
	amountBought := decimal.Zero
	acquiredCoin := decimal.Zero
	for _, level := range normalisedPrices {
		levelPrice := level.Price
		levelQuantity := level.Quantity

		if balance.LessThanOrEqual(levelQuantity) {
			quantity := balance
			balance = decimal.Zero
			amountBought = quantity.Mul(levelPrice)
		}
		if balance.GreaterThan(levelQuantity) {
			quantity := levelQuantity
			balance = balance.Sub(quantity)
			amountBought = quantity.Mul(levelPrice)
		}

		acquiredCoin = acquiredCoin.Add(amountBought)
		if balance.Equal(decimal.Zero) {
			return acquiredCoin
		}
	}
	return decimal.Zero

}
