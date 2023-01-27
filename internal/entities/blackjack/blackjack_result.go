package blackjack

import (
	"errors"
	internalErrors "scrub/internal/errors"

	"go.uber.org/zap"
)

var (
	win       = "win"
	loss      = "loss"
	push      = "push"
	blackjack = "blackjack"
)

func (bj *Blackjack) Results(logger *zap.Logger, players []BlackJackPlayer, dealerHand DealerHand) error {
	logger.Info("calculating results")

	err := bj.DrawDealerCards(logger, &dealerHand)
	if err != nil {
		return internalErrors.ErrFailedSubMethod("DrawDealerCards", err)
	}

	dealerBust := false
	if dealerHand.Bust() {
		logger.Info("dealer bust")
		dealerBust = true
	}

	for i, p := range players {
		for j, h := range p.Hands {
			if h.Bust() {
				players[i].Hands[j].result = &loss
				continue
			}

			if h.Blackjack() && !dealerHand.Blackjack() {
				players[i].Hands[j].result = &blackjack
				continue
			}

			if dealerBust {
				players[i].Hands[j].result = &win
				continue
			}

			if h.UpperValue() < dealerHand.UpperValue() {
				players[i].Hands[j].result = &loss
				continue
			}

			if h.UpperValue() > dealerHand.UpperValue() {
				players[i].Hands[j].result = &win
				continue
			}

			if h.UpperValue() == dealerHand.UpperValue() {
				players[i].Hands[j].result = &push
				continue
			}

			return errors.New("unexpected case")
		}
	}

	for i, p := range players {
		for _, h := range p.Hands {
			if h.result == nil {
				return errors.New("unexpected nil result")
			}
			res := *h.result

			switch res {
			case win:
				players[i].Win(h.betAmount)
			case loss:
				err = players[i].Lose(h.betAmount)
				if err != nil {
					return internalErrors.ErrFailedSubMethod("Lose", err)
				}
			case push:
				players[i].Win(0)
			case blackjack:
				players[i].Win(uint64(float64(h.betAmount) * 1.5))
			}
		}
	}

	return nil
}
