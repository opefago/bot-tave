package uniswap

import (
	"encoding/json"
	"fmt"
	"os"

	botcommon "github.com/opefago/bot-tave/common"
)

func New() *Uniswap {
	primaryTickers := []string{"USDT", "USDC"}
	uniswap := &Uniswap{
		BaseUrl: "https://uniswap.com",
	}
	uniswap.PrimaryTickers = primaryTickers
	return uniswap
}

// func getPiceData(baseUrl string) (map[string]Price, error) {
// 	url := fmt.Sprintf("%s/public?command=returnTicker", baseUrl)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	prices := make(map[string]Price)

// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(body, &prices)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return prices, nil
// }

func (uniswap *Uniswap) GetAssetPairs() ([]botcommon.TickerPair, error) {
	pairs := make([]botcommon.TickerPair, 100)

	return pairs, nil
}

func (uniswap *Uniswap) GetName() string {
	return "uniswap"
}

func (uniswap *Uniswap) StoreTriangularPair() error {
	prices, err := uniswap.GetAssetPairs()
	if err != nil {
		return err
	}

	fmt.Printf("Length of  pair %d\n", len(prices))

	triangularPair := botcommon.GetTraingularPairs(prices)

	fmt.Printf("Length of triangular pair %d\n", len(triangularPair))

	jsonOut, err := json.Marshal(triangularPair)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s.json", uniswap.GetName()))

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

func (uniswap *Uniswap) RunExchange() {}

func (uniswap *Uniswap) CalculateSurfaceArbitrage() []botcommon.TraingularPair {
	surfaceArb := make([]botcommon.TraingularPair, 5)

	return surfaceArb
}
