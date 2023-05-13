package blackjack

import (
	"errors"
	"scrub/internal/entities/player"
	"scrub/internal/errorutils"

	"go.uber.org/zap"
)

type BlackjackPlayer struct {
	player.Player
	Hands []Hand
}

func NewBlackjackPlayer(player player.Player, hands []Hand) BlackjackPlayer {
	return BlackjackPlayer{
		Player: player,
		Hands:  hands,
	}
}

func (bjp *BlackjackPlayer) PrintResult(logger *zap.Logger) error {
	for i, h := range bjp.Hands {
		if h.Result == nil {
			return errorutils.ErrUnexpectedNil
		}
		logger.Debug("player result",
			zap.String("player", bjp.Player.Name),
			zap.Int("hand", i+1),
			zap.Uint("hand value", h.UpperValue()),
			zap.String("result", *h.Result),
			zap.Uint64("credits", bjp.Player.Credits))
	}

	return nil
}

func PrintAllResults(logger *zap.Logger, blackjackPlayers []BlackjackPlayer) error {
	for i := range blackjackPlayers {
		if err := blackjackPlayers[i].PrintResult(logger); err != nil {
			return err
		}
	}
	return nil
}

func (bjp *BlackjackPlayer) ResetHands() error {
	splitCount := 0

	for i := range bjp.Hands {
		if bjp.Hands[i].isSplit {
			splitCount++
			bjp.Hands[i].isSplit = false
		}

		bjp.Hands[i].ResetDouble()
	}

	if splitCount%2 != 0 {
		return errors.New("split count is odd")
	}

	if splitCount > 0 {
		bjp.Hands = bjp.Hands[splitCount/2:]
	}

	return nil
}
