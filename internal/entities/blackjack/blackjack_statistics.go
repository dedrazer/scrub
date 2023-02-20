package blackjack

import "go.uber.org/zap"

func (bj *Blackjack) LogStatistics(logger *zap.Logger) {
	logger.Info("blackjack result statistics",
		zap.Uint64("player wins", bj.PlayerWins),
		zap.Uint64("player losses", bj.PlayerLosses),
		zap.Uint64("pushes", bj.Pushes),
		zap.Uint64("blackjacks", bj.PlayerBlackjackCount))

	logger.Info("blackjack additional statistics",
		zap.Uint64("splits", bj.SplitCount),
		zap.Uint64("dealer busts", bj.DealerBust),
		zap.Uint64("player busts", bj.PlayerBust))
}
