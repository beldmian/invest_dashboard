package investdata

import "context"

type InvestDataProvider interface {
	GetPortfolio(ctx context.Context) (*Portfolio, error)
}
