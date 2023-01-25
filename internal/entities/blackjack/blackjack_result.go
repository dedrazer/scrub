package blackjack

import (
	"scrub/internal/errors"

	"go.uber.org/zap"
)

var (
	win  = true
	loss = false
)

type Result struct {
	Status      *bool
	Credit      int
	IsBlackjack bool
}

func (bj *Blackjack) Results(logger *zap.Logger, playerHands []Hand, dealerHand DealerHand) ([]Result, error) {
	logger.Info("calculating results")
	results := make([]Result, len(playerHands))

	err := bj.DrawDealerCards(logger, &dealerHand)
	if err != nil {
		return nil, errors.ErrFailedSubMethod("DrawDealerCards", err)
	}

	dealerBust := false
	if dealerHand.Bust() {
		logger.Info("dealer bust")
		dealerBust = true
		for i := range results {
			// default all players who have not busted to win
			if !playerHands[i].Bust() {
				results[i].Status = &win
			}
		}
	}

	for i, h := range playerHands {
		if h.Bust() {
			results[i].Status = &loss
			continue
		}

		if !dealerBust {
			if h.UpperValue() < dealerHand.UpperValue() {
				results[i].Status = &loss
				continue
			}
		}

		if h.UpperValue() == 21 {
			results[i].IsBlackjack = true
		}

		if h.UpperValue() > dealerHand.UpperValue() {
			results[i].Status = &win
			continue
		}
	}

	return results, nil
}
