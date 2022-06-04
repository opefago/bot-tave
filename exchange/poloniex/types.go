package poloniex

import (
	"github.com/opefago/bot-tave/exchange/base"
	"github.com/shopspring/decimal"
)

type Price struct {
	Id            int32           `json:"id"`
	Last          string          `json:"last"`
	LowestAsk     decimal.Decimal `json:"lowestAsk"`
	HighestBid    decimal.Decimal `json:"highestBid"`
	PercentChange decimal.Decimal `json:"percentChange"`
	BaseVolume    decimal.Decimal `json:"baseVolume"`
	QuoteVolume   decimal.Decimal `json:"quoteVolume"`
	IsFrozen      string          `json:"isFrozen"`
	PostOnly      string          `json:"postOnly"`
	High24hr      decimal.Decimal `json:"high24hr"`
	Low24hr       decimal.Decimal `json:"low24hr"`
}

type PriceDepth struct {
	Asks     [][2]decimal.Decimal `json:"asks"`
	Bids     [][2]decimal.Decimal `json:"bids"`
	IsFrozen string               `json:"isFrozen"`
	PostOnly string               `json:"postOnly"`
	Seq      int                  `json:"seq"`
}

type Poloniex struct {
	BaseUrl string
	base.BaseExchange
}
