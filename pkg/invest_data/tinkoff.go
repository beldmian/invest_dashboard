package investdata

import (
	"context"

	"github.com/beldmian/invest_dashboard/pkg/config"
	"github.com/tinkoff/invest-api-go-sdk/investgo"
	investapi "github.com/tinkoff/invest-api-go-sdk/proto"
	"go.uber.org/zap"
)

type InvestDataProviderTinkoff struct {
	accountId                string
	operationsServiceClient  *investgo.OperationsServiceClient
	instrumentsServiceClient *investgo.InstrumentsServiceClient
	l                        *zap.Logger
}

func (i *InvestDataProviderTinkoff) GetPortfolio(ctx context.Context) (*Portfolio, error) {
	portfolioTinkoff, err := i.operationsServiceClient.GetPortfolio(i.accountId, investapi.PortfolioRequest_RUB)
	if err != nil {
		i.l.Error("GetPortfolio error", zap.Error(err))
		return nil, err
	}
	portfolio := Portfolio{
		BrokerID:   "tinkoff",
		TotalValue: portfolioTinkoff.GetTotalAmountPortfolio().ToFloat(),
		Positions:  make([]PortfolioPosition, len(portfolioTinkoff.GetPositions())),
	}
	for ind, pos := range portfolioTinkoff.GetPositions() {
		instrument, err := i.instrumentsServiceClient.InstrumentByFigi(pos.GetFigi())
		if err != nil {
			i.l.Error("GetInstrumentByFigi error", zap.Error(err))
			return nil, err
		}
		portfolio.Positions[ind] = PortfolioPosition{
			FigiID:         pos.GetFigi(),
			InstrumentType: pos.GetInstrumentType(),
			Name:           instrument.GetInstrument().GetName(),
			Ticker:         instrument.Instrument.GetTicker(),
			AveragePrice:   pos.GetAveragePositionPrice().ToFloat(),
			CurrentPrice:   pos.GetCurrentPrice().ToFloat(),
			Quantity:       pos.GetQuantity().ToFloat(),
		}
	}
	return &portfolio, nil
}

func ProvideInvestDataProviderTinkoff(c *config.Config, l *zap.Logger) *InvestDataProviderTinkoff {
	ctx := context.Background()
	logger := l.Sugar()
	client, err := investgo.NewClient(ctx, c.TinkoffInvestConfig, logger)
	if err != nil {
		l.Fatal("InvestDataProvider.NewClient error", zap.Error(err))
	}
	operationsServiceClient := client.NewOperationsServiceClient()
	instrumentsServiceClient := client.NewInstrumentsServiceClient()
	provider := InvestDataProviderTinkoff{
		accountId:                c.TinkoffInvestConfig.AccountId,
		operationsServiceClient:  operationsServiceClient,
		instrumentsServiceClient: instrumentsServiceClient,
		l:                        l,
	}
	return &provider
}
