package config

import "github.com/tinkoff/invest-api-go-sdk/investgo"

type Config struct {
	TinkoffInvestConfig investgo.Config
}

func ProvideConfig() *Config {
	investConfig, err := investgo.LoadConfig("tinkoff_config.yaml")
	if err != nil {
		panic(err)
	}
	return &Config{
		TinkoffInvestConfig: investConfig,
	}
}
