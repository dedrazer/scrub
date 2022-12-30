package blackjack

import (
	"scrub/internal/errors"

	"go.uber.org/zap"
)

var (
	win  = true
	loss = false
)

func (bj *Blackjack) Results(logger *zap.Logger, playerHands []Hand, dealerHand DealerHand) ([]*bool, error) {
	logger.Info("calculating results")
	results := make([]*bool, len(playerHands))

	err := bj.DrawDealerCards(logger, &dealerHand)
	if err != nil {
		return nil, errors.ErrFailedSubMethod("DrawDealerCards", err)
	}

	if dealerHand.Bust() {
		logger.Info("dealer bust")
		for i := range results {
			results[i] = &win
		}
		return results, nil
	}

	for i, h := range playerHands {
		if h.Bust() {
			results[i] = &loss
			continue
		}

		if h.UpperValue() < dealerHand.UpperValue() {
			results[i] = &loss
			continue
		}

		if h.UpperValue() > dealerHand.UpperValue() {
			results[i] = &win
			continue
		}
	}

	return results, nil
}
