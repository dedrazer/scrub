package main

import (
	"fmt"
	"scrub/internal/config"
	blackjackanalytics "scrub/internal/entities/analytics/blackjack"
	"scrub/internal/entities/blackjack/bettingstrategy"
	"scrub/internal/utils"

	"go.uber.org/zap"
)

func main() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableCaller = true
	zapConfig.Level = zap.NewAtomicLevelAt(utils.ParseZapLevel(config.LogLevel))

	logger, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to init zap: %s ", err.Error()))
	}

	simulationConfig := blackjackanalytics.SimulationConfig{
		MaxRounds:       200000,
		Decks:           6,
		BankCredits:     3000000,
		BankAtCredits:   10000,
		StartingCredits: 3000,
		OneCreditAmount: 50,
		RebuyCount:      20000,
	}

	strategyConfigs := bettingstrategy.CommonStrategyVariables{
		OneCreditValue: simulationConfig.OneCreditAmount,
		Logger:         logger,
	}

	strategies := []bettingstrategy.Strategy{
		&bettingstrategy.Stern{
			CommonStrategyVariables: strategyConfigs,
		},
		&bettingstrategy.Martingale{
			CommonStrategyVariables: strategyConfigs,
		},
	}

	for i := range strategies {
		simulator := blackjackanalytics.NewSimulator(logger, strategies[i], simulationConfig)

		err = simulator.Simulate()
		if err != nil {
			logger.Fatal("unexpected error", zap.Error(err))
		}
	}
}
