package updater

import (
	"context"
	"time"

	investdata "github.com/beldmian/invest_dashboard/pkg/invest_data"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PortfolioMetrics struct {
	TotalPortfolioPrice   prometheus.Gauge
	PositionsCurrentPrice map[string]prometheus.Gauge
	PositionsAveragePrice map[string]prometheus.Gauge
	PositionsQuantity     map[string]prometheus.Gauge
}

type Updater struct {
	investDataProvider investdata.InvestDataProvider
	metrics            *PortfolioMetrics
	l                  *zap.Logger
}

func (u *Updater) Start() {
	for {
		portfolio, err := u.investDataProvider.GetPortfolio(context.Background())
		if err != nil {
			u.l.Warn("GetPortfolio error", zap.Error(err))
			continue
		}
		u.l.Info("Got portfolio info", zap.Time("now", time.Now()))
		u.metrics.TotalPortfolioPrice.Set(portfolio.TotalValue)
		for _, pos := range portfolio.Positions {
			if _, ok := u.metrics.PositionsCurrentPrice[pos.FigiID]; !ok {
				continue
			}
			u.metrics.PositionsCurrentPrice[pos.FigiID].Set(pos.CurrentPrice)
			u.metrics.PositionsAveragePrice[pos.FigiID].Set(pos.AveragePrice)
			u.metrics.PositionsQuantity[pos.FigiID].Set(pos.Quantity)
		}
		time.Sleep(15 * time.Second)
	}
}

func ProvideUpdater(investDataProvider investdata.InvestDataProvider, l *zap.Logger, lc fx.Lifecycle) *Updater {
	initialPortfolioData, err := investDataProvider.GetPortfolio(context.Background())
	if err != nil {
		panic(err)
	}
	positionsCurrentPrice := make(map[string]prometheus.Gauge)
	positionsAveragePrice := make(map[string]prometheus.Gauge)
	positionsQuantity := make(map[string]prometheus.Gauge)
	for _, pos := range initialPortfolioData.Positions {
		positionsCurrentPrice[pos.FigiID] = promauto.NewGauge(prometheus.GaugeOpts{
			Subsystem:   initialPortfolioData.BrokerID,
			Name:        "tikcer_data_" + pos.Ticker + "_current",
			Help:        pos.Name,
			ConstLabels: prometheus.Labels{"type": "ticker_current", "ticker": pos.Ticker, "name": pos.Name},
		})
		positionsAveragePrice[pos.FigiID] = promauto.NewGauge(prometheus.GaugeOpts{
			Subsystem:   initialPortfolioData.BrokerID,
			Name:        "tikcer_data_" + pos.Ticker + "_average",
			Help:        pos.Name,
			ConstLabels: prometheus.Labels{"type": "ticker_average", "ticker": pos.Ticker, "name": pos.Name},
		})
		positionsQuantity[pos.FigiID] = promauto.NewGauge(prometheus.GaugeOpts{
			Subsystem:   initialPortfolioData.BrokerID,
			Name:        "tikcer_data_" + pos.Ticker + "_quantity",
			Help:        pos.Name,
			ConstLabels: prometheus.Labels{"type": "ticker_quantity", "ticker": pos.Ticker, "name": pos.Name},
		})
	}
	u := Updater{
		investDataProvider: investDataProvider,
		metrics: &PortfolioMetrics{
			TotalPortfolioPrice: promauto.NewGauge(prometheus.GaugeOpts{
				Subsystem: initialPortfolioData.BrokerID,
				Name:      "total_portfolio_price",
				Help:      "Total portfolio price",
			}),
			PositionsCurrentPrice: positionsCurrentPrice,
			PositionsAveragePrice: positionsAveragePrice,
			PositionsQuantity:     positionsQuantity,
		},
		l: l,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go u.Start()
			return nil
		},
	})
	return &u
}
