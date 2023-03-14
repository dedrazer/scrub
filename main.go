package main

import (
	"fmt"
	blackjackanalytics "scrub/internal/entities/analytics/blackjack"
	"scrub/internal/entities/blackjack/bettingstrategy"

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
	simulationConfig := blackjackanalytics.SimulationConfig{
		Rounds:          1000000,
		Decks:           6,
		BankCredits:     3000000,
		BankAtCredits:   10000,
		StartingCredits: 3000,
		OneCreditAmount: 50,
		RebuyCount:      20000,
	}

	err = blackjackanalytics.Simulate(logger, simulationConfig, bettingstrategy.Stern{})
	if err != nil {
		logger.Fatal("unexpected error", zap.Error(err))
	}
}
