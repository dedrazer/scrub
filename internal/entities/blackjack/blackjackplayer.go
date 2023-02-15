package blackjack

import (
	"scrub/internal/entities/player"
	"scrub/internal/errors"

	"go.uber.org/zap"
)

type BlackJackPlayer struct {
	player.Player
	Hands []Hand
}

func (bjp *BlackJackPlayer) PrintResult(logger *zap.Logger) error {
	for i, h := range bjp.Hands {
		if h.result == nil {
			return errors.ErrUnexpectedNil
		}
		logger.Debug("player result",
			zap.String("player", bjp.Player.Name),
			zap.Int("hand", i+1),
			zap.Uint("hand value", h.UpperValue()),
			zap.String("result", *h.result),
			zap.Uint64("credits", bjp.Player.Credits))
	}

	return nil
}

func PrintAllResults(logger *zap.Logger, blackjackPlayers []BlackJackPlayer) error {
	for i := range blackjackPlayers {
		if err := blackjackPlayers[i].PrintResult(logger); err != nil {
			return err
		}
	}
	return nil
}
