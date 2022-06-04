package base

import (
	"context"

	botcommon "github.com/opefago/bot-tave/common"
	"github.com/shopspring/decimal"
)

type Exchange interface {
	GetAssetPairs() ([]botcommon.TickerPair, error)
	GetName() string
	FetchExchangePrices(triangularPairs []botcommon.TraingularPair) ([]botcommon.TraingularPairPriceInfo, error)
	StoreTriangularPair() error
	CalculateSurfaceArbitrage() ([]botcommon.SurfaceRatePair, error)
	RunExchange(ctx context.Context)
	CalculateRealArbitrage(surfaceArb []botcommon.SurfaceRatePair) ([]botcommon.RealRatePair, error)
}

type BaseExchange struct {
	// PotentialPair []botcommon.TraingularPair
	// PotentialPair []botcommon.TraingularPair
	PrimaryTickers []string
	StartAmount    decimal.Decimal
}
