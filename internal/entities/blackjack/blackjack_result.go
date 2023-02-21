package blackjack

import (
	"errors"
	"scrub/internal/entities/blackjack/utils"
	internalErrors "scrub/internal/errors"

	"go.uber.org/zap"
)

func (bj *Blackjack) Results(logger *zap.Logger, players []BlackjackPlayer, dealerHand DealerHand) error {
	logger.Debug("calculating results")

	err := bj.DrawDealerCards(logger, &dealerHand)
	if err != nil {
		return internalErrors.ErrFailedSubMethod("DrawDealerCards", err)
	}

	dealerBust := false
	if dealerHand.Bust() {
		logger.Debug("dealer bust")
		dealerBust = true
	}

	for i, p := range players {
		for j, h := range p.Hands {
			if h.Bust() {
				bj.PlayerBust++
				bj.PlayerLosses++
				players[i].Hands[j].Result = &utils.Loss
				continue
			}

			if h.Blackjack() && !dealerHand.Blackjack() {
				bj.PlayerBlackjackCount++
				bj.PlayerWins++
				players[i].Hands[j].Result = &utils.Blackjack
				continue
			}

			if dealerBust {
				bj.DealerBust++
				bj.PlayerWins++
				players[i].Hands[j].Result = &utils.Win
				continue
			}

			if h.UpperValue() < dealerHand.UpperValue() {
				bj.PlayerLosses++
				players[i].Hands[j].Result = &utils.Loss
				continue
			}

			if h.UpperValue() > dealerHand.UpperValue() {
				bj.PlayerWins++
				players[i].Hands[j].Result = &utils.Win
				continue
			}

			if h.UpperValue() == dealerHand.UpperValue() {
				bj.Pushes++
				players[i].Hands[j].Result = &utils.Push
				continue
			}

			return errors.New("unexpected case")
		}
	}

	for i, p := range players {
		for _, h := range p.Hands {
			if h.Result == nil {
				return errors.New("unexpected nil result")
			}
			res := *h.Result

			switch res {
			case utils.Win:
				players[i].Win(h.BetAmount)
			case utils.Loss:
				err = players[i].Lose(h.BetAmount)
				if err != nil {
					return internalErrors.ErrFailedSubMethod("Lose", err)
				}
			case utils.Push:
				players[i].Win(0)
			case utils.Blackjack:
				players[i].Win(uint64(float64(h.BetAmount) * 1.5))
			}
		}
	}

	return nil
}
