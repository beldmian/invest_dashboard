package investdata

type Portfolio struct {
	BrokerID   string
	TotalValue float64
	Positions  []PortfolioPosition
}

type PortfolioPosition struct {
	FigiID         string
	InstrumentType string
	Name           string
	Ticker         string
	AveragePrice   float64
	CurrentPrice   float64
	Quantity       float64
}
