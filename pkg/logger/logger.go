package logger

import "go.uber.org/zap"

func ProvideLogger() *zap.Logger {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return l
}
