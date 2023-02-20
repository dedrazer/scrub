package main

import (
	"fmt"
	blackjackanalytics "scrub/internal/entities/analytics/blackjack"

	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to init zap: %s ", err.Error()))
	}

	//deck.Demo(logger)
	//blackjack.Demo(logger)
	err = blackjackanalytics.Simulate(logger, 10000, 6)
	if err != nil {
		logger.Fatal("unexpected error", zap.Error(err))
	}
}
