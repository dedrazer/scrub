package blackjack

import (
	"scrub/internal/entities/player"
	"scrub/internal/errors"

	"go.uber.org/zap"
)

type BlackJackPlayer struct {
	PlayerBet player.PlayerBet
	Hands     []Hand
}

func (bjp BlackJackPlayer) PrintResult(logger *zap.Logger) error {
	for i, h := range bjp.Hands {
		if h.result == nil {
			return errors.ErrUnexcpectedNil
		}
		logger.Info("player result", zap.Int("hand", i+1), zap.String("result", *h.result))
	}
	
	return nil
}
