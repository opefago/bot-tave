package binance

import "github.com/opefago/bot-tave/exchange/base"

type Binance struct {
	BaseUrl string
	base.BaseExchange
}

type exchangeInfo struct {
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol                     string   `json:"symbol"`
		Status                     string   `json:"status"`
		BaseAsset                  string   `json:"baseAsset"`
		BaseAssetPrecision         int      `json:"baseAssetPrecision"`
		QuoteAsset                 string   `json:"quoteAsset"`
		QuotePrecision             int      `json:"quotePrecision"`
		QuoteAssetPrecision        int      `json:"quoteAssetPrecision"`
		BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
		QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
		OrderTypes                 []string `json:"orderTypes"`
		IcebergAllowed             bool     `json:"icebergAllowed"`
		OcoAllowed                 bool     `json:"ocoAllowed"`
		QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
		AllowTrailingStop          bool     `json:"allowTrailingStop"`
		IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
		IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
		Filters                    []struct {
			FilterType            string `json:"filterType"`
			MinPrice              string `json:"minPrice,omitempty"`
			MaxPrice              string `json:"maxPrice,omitempty"`
			TickSize              string `json:"tickSize,omitempty"`
			MultiplierUp          string `json:"multiplierUp,omitempty"`
			MultiplierDown        string `json:"multiplierDown,omitempty"`
			AvgPriceMins          int    `json:"avgPriceMins,omitempty"`
			MinQty                string `json:"minQty,omitempty"`
			MaxQty                string `json:"maxQty,omitempty"`
			StepSize              string `json:"stepSize,omitempty"`
			MinNotional           string `json:"minNotional,omitempty"`
			ApplyToMarket         bool   `json:"applyToMarket,omitempty"`
			Limit                 int    `json:"limit,omitempty"`
			MinTrailingAboveDelta int    `json:"minTrailingAboveDelta,omitempty"`
			MaxTrailingAboveDelta int    `json:"maxTrailingAboveDelta,omitempty"`
			MinTrailingBelowDelta int    `json:"minTrailingBelowDelta,omitempty"`
			MaxTrailingBelowDelta int    `json:"maxTrailingBelowDelta,omitempty"`
			MaxNumOrders          int    `json:"maxNumOrders,omitempty"`
			MaxNumAlgoOrders      int    `json:"maxNumAlgoOrders,omitempty"`
		} `json:"filters"`
		Permissions []string `json:"permissions"`
	} `json:"symbols"`
}
