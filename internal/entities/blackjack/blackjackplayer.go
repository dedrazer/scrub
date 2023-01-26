package blackjack

import (
	"scrub/internal/entities/player"

	"go.uber.org/zap"
)

type BlackJackPlayer struct {
	PlayerBet player.PlayerBet
	Hands     []Hand
}

func (bjp BlackJackPlayer) PrintResult(logger *zap.Logger) {
	bjp.PrintResult(logger)
}
