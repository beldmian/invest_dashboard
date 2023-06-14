package main

import (
	"github.com/beldmian/invest_dashboard/internal/pkg/metrics"
	"github.com/beldmian/invest_dashboard/internal/pkg/updater"
	"github.com/beldmian/invest_dashboard/pkg/config"
	investdata "github.com/beldmian/invest_dashboard/pkg/invest_data"
	"github.com/beldmian/invest_dashboard/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.ProvideConfig,
			logger.ProvideLogger,
		),
		fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: l}
		}),
		fx.Provide(
			fx.Annotate(
				investdata.ProvideInvestDataProviderTinkoff,
				fx.As(new(investdata.InvestDataProvider)),
			),
		),
		fx.Provide(updater.ProvideUpdater),
		fx.Invoke(
			func(u *updater.Updater) {},
			metrics.InvokeMetricsServer,
		),
	)
	app.Run()
}
