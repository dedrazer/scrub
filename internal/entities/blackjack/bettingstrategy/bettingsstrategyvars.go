package bettingstrategy

import "go.uber.org/zap"

type CommonStrategyVariables struct {
	Logger         *zap.Logger
	OneCreditValue uint64
	lossStreak     int
	winStreak      int
	round          int
}
