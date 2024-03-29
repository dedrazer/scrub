package bettingstrategy

import (
	"scrub/internal/entities/blackjack"

	"go.uber.org/zap"
)

type Strategy interface {
	GetName() string
	BettingStrategy(players []blackjack.BlackjackPlayer) error
}

func playerAllIn(logger *zap.Logger, player *blackjack.BlackjackPlayer, handNumber, round, lossStreak int, roundResult string) {
	player.Hands[handNumber].BetAmount = player.Credits
	logger.Debug("player all in",
		zap.String("round result", roundResult),
		zap.Int("loss streak", lossStreak),
		zap.Uint64("next bet", player.Hands[handNumber].BetAmount),
		zap.Int("round", round))
}
