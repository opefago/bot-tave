package poloniex

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	botcommon "github.com/opefago/bot-tave/common"
	"github.com/opefago/bot-tave/utils"
	"github.com/shopspring/decimal"
)

func New() *Poloniex {
	primaryTickers := []string{"USDT", "USDC", "USDD", "USDJ", "TUSD"}
	amount, _ := decimal.NewFromString("100")
	poloniex := &Poloniex{
		BaseUrl: "https://poloniex.com",
	}
	poloniex.PrimaryTickers = primaryTickers
	poloniex.StartAmount = amount
	return poloniex
}

func getPiceData(baseUrl string) (map[string]Price, error) {
	url := fmt.Sprintf("%s/public?command=returnTicker", baseUrl)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	prices := make(map[string]Price)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &prices)

	if err != nil {
		return nil, err
	}

	return prices, nil
}

func (pol *Poloniex) GetAssetPairs() ([]botcommon.TickerPair, error) {
	pairs := make([]botcommon.TickerPair, 100)

	prices, err := getPiceData(pol.BaseUrl)

	if err != nil {
		return nil, err
	}

	for pair, _ := range prices {
		splitPair := strings.Split(pair, "_")
		pair := botcommon.TickerPair{
			Base:  splitPair[0],
			Quote: splitPair[1],
		}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

func (pol *Poloniex) GetName() string {
	return "poloniex"
}

func (poloniex *Poloniex) StoreTriangularPair() error {
	prices, err := poloniex.GetAssetPairs()
	if err != nil {
		return err
	}

	fmt.Printf("Length of  pair %d\n", len(prices))

	triangularPair := botcommon.GetTraingularPairs(prices)
	jsonOut, err := json.Marshal(triangularPair)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s.json", poloniex.GetName()))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(string(jsonOut))

	if err != nil {
		return err
	}

	sortedTriangularPair := utils.FilterByStableCoins(triangularPair, poloniex.PrimaryTickers)
	jsonOut, err = json.Marshal(sortedTriangularPair)

	if err != nil {
		return err
	}

	f_stable, err := os.Create(fmt.Sprintf("%s_stable.json", poloniex.GetName()))
	if err != nil {
		return err
	}
	defer f_stable.Close()

	_, err = f_stable.WriteString(string(jsonOut))

	if err != nil {
		return err
	}

	return nil
}

func (poloniex *Poloniex) ExecuteTrade(realArb botcommon.RealRatePair) error {
	return nil
}

func (poloniex *Poloniex) RunExchange(ctx context.Context) {
	t := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-t.C:
			surfaceArb, err := poloniex.CalculateSurfaceArbitrage()
			if err != nil {
				continue
			}
			if len(surfaceArb) > 0 {
				realArb, err := poloniex.CalculateRealArbitrage(surfaceArb)
				if err != nil {
					continue
				}
				for _, arb := range realArb {
					poloniex.ExecuteTrade(arb)
				}
			}
		case <-ctx.Done():
			return
		}

	}
}

func getPriceDepth(baseUrl string, tickerPair botcommon.TickerPair) (PriceDepth, error) {
	combined := fmt.Sprintf("%s_%s", tickerPair.Base, tickerPair.Quote)
	url := fmt.Sprintf("%s/public?command=returnOrderBook&currencyPair=%s&depth=10", baseUrl, combined)
	resp, err := http.Get(url)
	if err != nil {
		return PriceDepth{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var priceDepth PriceDepth

	if err != nil {
		return PriceDepth{}, err
	}

	err = json.Unmarshal(body, &priceDepth)

	if err != nil {
		return PriceDepth{}, err
	}

	return priceDepth, nil
}

func normalizeOrderBook(priceDepth PriceDepth, direction botcommon.Direction) []botcommon.PriceQuantity {
	var normailizedPrices []botcommon.PriceQuantity
	switch direction {
	case botcommon.Forward:
		for _, price := range priceDepth.Asks {
			askPrice := decimal.NewFromInt(1).Div(price[0])
			adjQty := price[1].Mul(askPrice)
			normailizedPrices = append(normailizedPrices, botcommon.PriceQuantity{
				Price:    askPrice,
				Quantity: adjQty,
			})
		}
	case botcommon.Reverse:
		for _, price := range priceDepth.Bids {
			bidPrice := price[0]
			adjQty := price[1]
			normailizedPrices = append(normailizedPrices, botcommon.PriceQuantity{
				Price:    bidPrice,
				Quantity: adjQty,
			})
		}
	}
	return normailizedPrices
}

func (poloniex *Poloniex) CalculateRealArbitrage(surfaceArb []botcommon.SurfaceRatePair) ([]botcommon.RealRatePair, error) {
	var realArb []botcommon.RealRatePair
	for _, rate := range surfaceArb {
		primaryDepth, err := getPriceDepth(poloniex.BaseUrl, rate.PrimaryRatePair.Pair)
		if err != nil {
			return nil, err
		}
		primaryNomalized := normalizeOrderBook(primaryDepth, rate.PrimaryRatePair.Direction)

		secondaryDepth, err := getPriceDepth(poloniex.BaseUrl, rate.SecondaryRatePair.Pair)
		if err != nil {
			return nil, err
		}
		secondaryNomalized := normalizeOrderBook(secondaryDepth, rate.SecondaryRatePair.Direction)

		tertiaryDepth, err := getPriceDepth(poloniex.BaseUrl, rate.TertiayRatePair.Pair)
		if err != nil {
			return nil, err
		}
		tertiaryNomalized := normalizeOrderBook(tertiaryDepth, rate.TertiayRatePair.Direction)

		primaryResult := utils.CalculateRealPotential(poloniex.StartAmount, primaryNomalized)
		secondaryResult := utils.CalculateRealPotential(primaryResult, secondaryNomalized)
		tertiaryResult := utils.CalculateRealPotential(secondaryResult, tertiaryNomalized)

		pnl := tertiaryResult.Sub(poloniex.StartAmount)

		if pnl.GreaterThan(decimal.Zero) {
			arb := botcommon.RealRatePair{}
			arb.SurfaceRatePair = rate
			arb.RealProfit = pnl
			realArb = append(realArb, arb)
		}

		fmt.Printf("Tertiary Result %s\n", tertiaryResult)
		fmt.Printf("Real Profit/Loss %s\n", pnl.Div(poloniex.StartAmount))
		fmt.Println("===========================================================")
		fmt.Println()
		fmt.Println()

	}
	return realArb, nil
}

func (poloniex *Poloniex) CalculateSurfaceArbitrage() ([]botcommon.SurfaceRatePair, error) {
	var surfaceArb []botcommon.SurfaceRatePair

	triangularPair, err := botcommon.LoadTriangularPair(fmt.Sprintf("%s_stable.json", poloniex.GetName()))
	if err != nil {
		return nil, err
	}
	prices, err := poloniex.FetchExchangePrices(triangularPair)
	if err != nil {
		return nil, err
	}

	for _, price := range prices {
		for _, stableCoin := range poloniex.PrimaryTickers {
			//checks if the coin starts and terminates pn a stable coin
			//transactions are only transacted for transactions starting and ending in a stable coin
			//This is to avoid holding a crypto asset for too long thereby exposing myself to crypto volatility
			if (price.PrimaryTickerPair.Pair.Base == stableCoin || price.PrimaryTickerPair.Pair.Quote == stableCoin) &&
				(price.TertiaryTickerPair.Pair.Base == stableCoin || price.TertiaryTickerPair.Pair.Quote == stableCoin) {

				startingAmount := decimal.NewFromInt(1)
				primaryAcquiredValue, primaryAcquiredToken, primaryDirection := utils.CalculateTokenReturns(stableCoin, startingAmount, price.PrimaryTickerPair)

				secondaryAcquiredValue, secondaryAcquiredToken, secondaryDirection := utils.CalculateTokenReturns(primaryAcquiredToken, primaryAcquiredValue, price.SecondaryTickerPair)

				tertiaryAcquiredValue, _, tertiaryDirection := utils.CalculateTokenReturns(secondaryAcquiredToken, secondaryAcquiredValue, price.TertiaryTickerPair)

				pnl := tertiaryAcquiredValue.Sub(startingAmount)

				if pnl.GreaterThan(decimal.Zero) {
					// fmt.Printf("Starting with pricing information %+v\n", price)
					// fmt.Println()
					// fmt.Printf("Starting with token %s and amount %s: get token %s for %s\n", stableCoin, startingAmount, primaryAcquiredToken, primaryAcquiredValue)
					// fmt.Printf("Swaping token %s and amount %s: for token %s at %s\n", primaryAcquiredToken, primaryAcquiredValue, secondaryAcquiredToken, secondaryAcquiredValue)
					// fmt.Printf("Swaping token %s and amount %s: for token %s at %s\n", secondaryAcquiredToken, secondaryAcquiredValue, tertiaryAcquiredToken, tertiaryAcquiredValue)
					// fmt.Printf("result %s, percentage PNL %s\n", pnl, pnl.Div(startingAmount).Mul(decimal.NewFromInt(100)))

					primarySwap := botcommon.RatePair{
						Pair:      price.PrimaryTickerPair.Pair,
						Direction: primaryDirection,
						ValueIn:   startingAmount,
						ValueOut:  primaryAcquiredValue,
					}

					secondarySwap := botcommon.RatePair{
						Pair:      price.SecondaryTickerPair.Pair,
						Direction: secondaryDirection,
						ValueIn:   primaryAcquiredValue,
						ValueOut:  secondaryAcquiredValue,
					}
					tertiarySwap := botcommon.RatePair{
						Pair:      price.TertiaryTickerPair.Pair,
						Direction: tertiaryDirection,
						ValueIn:   secondaryAcquiredValue,
						ValueOut:  tertiaryAcquiredValue,
					}

					surfaceRate := botcommon.SurfaceRatePair{
						PrimaryRatePair:   primarySwap,
						SecondaryRatePair: secondarySwap,
						TertiayRatePair:   tertiarySwap,
						Profit:            pnl,
					}
					surfaceArb = append(surfaceArb, surfaceRate)
				}
			}
		}
	}

	return surfaceArb, nil
}

func (polonies *Poloniex) FetchExchangePrices(triangularPairs []botcommon.TraingularPair) ([]botcommon.TraingularPairPriceInfo, error) {
	var pricesInfo []botcommon.TraingularPairPriceInfo
	prices, err := getPiceData(polonies.BaseUrl)
	if err != nil {
		return nil, err
	}

	for _, triangularPair := range triangularPairs {
		var traingularPriceInfo botcommon.TraingularPairPriceInfo
		primary := getAdjustednfoWithPrice(triangularPair.PrimaryTickerPair, prices)
		traingularPriceInfo.PrimaryTickerPair = primary

		secondary := getAdjustednfoWithPrice(triangularPair.SecondaryTickerPair, prices)
		traingularPriceInfo.SecondaryTickerPair = secondary

		tertiary := getAdjustednfoWithPrice(triangularPair.TertiaryTickerPair, prices)
		traingularPriceInfo.TertiaryTickerPair = tertiary

		pricesInfo = append(pricesInfo, traingularPriceInfo)

	}

	return pricesInfo, nil
}

func getAdjustednfoWithPrice(tickerPair botcommon.TickerPair, prices map[string]Price) botcommon.PriceInfo {
	var priceInfo botcommon.PriceInfo
	pair := fmt.Sprintf("%s_%s", tickerPair.Base, tickerPair.Quote)
	price := prices[pair]
	priceInfo.Pair = botcommon.TickerPair{
		Base:  tickerPair.Base,
		Quote: tickerPair.Quote,
	}

	priceInfo.Price = botcommon.TickerPairPrice{
		Bid: price.HighestBid,
		Ask: price.LowestAsk,
	}
	return priceInfo
}
