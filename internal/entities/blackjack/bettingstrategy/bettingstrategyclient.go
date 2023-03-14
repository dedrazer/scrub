package bettingstrategy

import (
	"scrub/internal/entities/blackjack"

	"go.uber.org/zap"
)

type Strategy interface {
	GetName() string
	Strategy(players []blackjack.BlackjackPlayer) error
}

func playerAllIn(logger *zap.Logger, player *blackjack.BlackjackPlayer, handNumber, round, lossStreak int) {
	player.Hands[handNumber].BetAmount = player.Credits
	logger.Debug("player all in",
		zap.Int("loss streak", lossStreak),
		zap.Uint64("next bet", player.Hands[handNumber].BetAmount),
		zap.Int("round", round))
}
