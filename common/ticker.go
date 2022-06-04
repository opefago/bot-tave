package common

import (
	"github.com/shopspring/decimal"
)

type TickerPair struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

type RatePair struct {
	Pair      TickerPair      `json:"pair"`
	ValueIn   decimal.Decimal `json:"valueIn"`
	ValueOut  decimal.Decimal `json:"valueOut"`
	Direction Direction       `json:"direction"`
}

type SurfaceRatePair struct {
	PrimaryRatePair   RatePair
	SecondaryRatePair RatePair
	TertiayRatePair   RatePair
	Profit            decimal.Decimal
}

type RealRatePair struct {
	SurfaceRatePair
	RealProfit decimal.Decimal
}

type TickerPairPrice struct {
	Bid decimal.Decimal `json:"bid"`
	Ask decimal.Decimal `json:"ask"`
}

type PriceInfo struct {
	Pair  TickerPair
	Price TickerPairPrice
}

// func (pair TickerPair) String() string {
// 	return fmt.Sprintf("%s_%s", pair.Base, pair.Quote)
// }

type Direction int8

const (
	Forward Direction = iota
	Reverse
)

type TraingularPair struct {
	PrimaryTickerPair   TickerPair `json:"primaryTickerPair"`
	SecondaryTickerPair TickerPair `json:"secondaryTickerPair"`
	TertiaryTickerPair  TickerPair `json:"tertiaryTickerPair"`
}

type TraingularPairPriceInfo struct {
	PrimaryTickerPair   PriceInfo `json:"primaryTickerPair"`
	SecondaryTickerPair PriceInfo `json:"secondaryTickerPair"`
	TertiaryTickerPair  PriceInfo `json:"tertiaryTickerPair"`
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type PriceQuantity struct {
	Price    decimal.Decimal
	Quantity decimal.Decimal
}

func GetTraingularPairs(pairs []TickerPair) []TraingularPair {
	var triangularPair []TraingularPair
	length := len(pairs)
	for i := 0; i < length; i++ {
		primaryPair := []string{pairs[i].Base, pairs[i].Quote}
		for j := i + 1; j < length; j++ {
			secondaryPair := []string{pairs[j].Base, pairs[j].Quote}
			if Contains(primaryPair, secondaryPair[0]) || Contains(primaryPair, secondaryPair[1]) {
				for k := j + 1; k < length; k++ {
					tertiaryPair := []string{pairs[k].Base, pairs[k].Quote}
					if ((Contains(secondaryPair, tertiaryPair[0]) && !Contains(primaryPair, tertiaryPair[0])) ||
						(Contains(primaryPair, tertiaryPair[0]) && !Contains(secondaryPair, tertiaryPair[0]))) &&
						((Contains(secondaryPair, tertiaryPair[1]) && !Contains(primaryPair, tertiaryPair[1])) ||
							(Contains(primaryPair, tertiaryPair[1]) && !Contains(secondaryPair, tertiaryPair[1]))) {

						pair := TraingularPair{
							PrimaryTickerPair:   TickerPair{Base: primaryPair[0], Quote: primaryPair[1]},
							SecondaryTickerPair: TickerPair{Base: secondaryPair[0], Quote: secondaryPair[1]},
							TertiaryTickerPair:  TickerPair{Base: tertiaryPair[0], Quote: tertiaryPair[1]},
						}

						triangularPair = append(triangularPair, pair)
					}
				}
			}
		}
	}
	return triangularPair
}
