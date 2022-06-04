package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	botcommon "github.com/opefago/bot-tave/common"
	"github.com/opefago/bot-tave/utils"
)

func New() *Binance {
	primaryTickers := []string{"USDT", "BUSD", "USDC", "TUSD", "USDP"}
	binance := &Binance{BaseUrl: "https://api.binance.com/api/v3"}
	binance.PrimaryTickers = primaryTickers
	return binance
}

func (binance *Binance) getPiceData() {}

func (binance *Binance) GetAssetPairs() ([]botcommon.TickerPair, error) {
	url := fmt.Sprintf("%s/exchangeInfo", binance.BaseUrl)
	pairs := make([]botcommon.TickerPair, 3)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var info exchangeInfo

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &info)

	if err != nil {
		return nil, err
	}

	for _, pair := range info.Symbols {
		pair := botcommon.TickerPair{
			Base:  pair.BaseAsset,
			Quote: pair.QuoteAsset,
		}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

func (binance *Binance) GetName() string {
	return "binance"
}

func (binance *Binance) StoreTriangularPair() error {
	prices, err := binance.GetAssetPairs()
	if err != nil {
		return err
	}

	fmt.Printf("Length of  pair %d\n", len(prices))

	triangularPair := botcommon.GetTraingularPairs(prices)

	sortedTriangularPair := utils.FilterByStableCoins(triangularPair, binance.PrimaryTickers)

	fmt.Printf("Length of triangular pair %d\n", len(triangularPair))
	fmt.Printf("Length of Filter triangular pair %d\n", len(sortedTriangularPair))

	jsonOut, err := json.Marshal(triangularPair)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s.json", binance.GetName()))

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(string(jsonOut))

	if err != nil {
		return err
	}
	return nil
}

func (binance *Binance) RunExchange(ctx context.Context) {}

func (binance *Binance) CalculateSurfaceArbitrage() ([]botcommon.SurfaceRatePair, error) {
	surfaceArb := make([]botcommon.SurfaceRatePair, 5)
	return surfaceArb, nil
}

func (binance *Binance) CalculateRealArbitrage(surfaceArb []botcommon.SurfaceRatePair) ([]botcommon.RealRatePair, error) {
	return nil, nil
}

func (binance *Binance) FetchExchangePrices(triangularPairs []botcommon.TraingularPair) ([]botcommon.TraingularPairPriceInfo, error) {
	var pricesInfo []botcommon.TraingularPairPriceInfo
	return pricesInfo, nil
}

func (binance Binance) LoadTriangularPair() ([]botcommon.TraingularPair, error) {
	var triangularPairs []botcommon.TraingularPair
	return triangularPairs, nil
}
