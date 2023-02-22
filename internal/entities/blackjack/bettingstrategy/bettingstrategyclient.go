package bettingstrategy

import (
	"scrub/internal/entities/blackjack"

	"go.uber.org/zap"
)

func playerAllIn(logger *zap.Logger, player *blackjack.BlackjackPlayer, handNumber int) {
	player.Hands[handNumber].BetAmount = player.Credits
	logger.Debug("player all in",
		zap.Int("loss streak", lossStreak),
		zap.Uint64("next bet", player.Hands[handNumber].BetAmount),
		zap.Int("round", round))
}
