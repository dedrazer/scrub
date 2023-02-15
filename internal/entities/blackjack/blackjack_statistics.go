package blackjack

import "go.uber.org/zap"

func (bj *Blackjack) LogStatistics(logger *zap.Logger) {
	logger.Info("split counter", zap.Uint64("counter", bj.SplitCounter))
}
