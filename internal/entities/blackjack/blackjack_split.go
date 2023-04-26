package blackjack

import (
	"scrub/internal/entities/deck"
	"scrub/internal/errorutils"

	"go.uber.org/zap"
)

func (bj *Blackjack) split(player *BlackjackPlayer, handIndex int) error {
	if player.Hands[handIndex].cards[0].Symbol != player.Hands[handIndex].cards[1].Symbol {
		return errorutils.ErrCannotSplit
	}

	bj.logger.Debug("splitting hand", zap.Any("player", player), zap.Int("hand", handIndex+1))

	splitHand := Hand{
		cards:     []deck.Card{player.Hands[handIndex].cards[1]},
		isSplit:   true,
		BetAmount: player.Hands[handIndex].BetAmount,
	}
	player.Hands[handIndex].cards = player.Hands[handIndex].cards[:1]
	player.Hands[handIndex].isSplit = true
	player.Hands = append(player.Hands, splitHand)

	bj.SplitCount++

	return nil
}
