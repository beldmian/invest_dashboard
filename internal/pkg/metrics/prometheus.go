package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

func InvokeMetricsServer(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				http.Handle("/metrics", promhttp.Handler())
				http.ListenAndServe(":2112", nil)
			}()
			return nil
		},
	})
}
