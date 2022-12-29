package blackjack

import (
	"scrub/internal/errors"

	"go.uber.org/zap"
)

var t = true
var f = false

func (bj *Blackjack) Results(logger *zap.Logger, playerHands []Hand, dealerHand Hand) ([]*bool, error) {
	logger.Info("calculating results")
	results := make([]*bool, len(playerHands))

	err := bj.DrawDealerCards(logger, &dealerHand)
	if err != nil {
		return nil, errors.ErrFailedSubMethod("DrawDealerCards", err)
	}

	if dealerHand.Bust() {
		logger.Info("dealer bust")
		for i := range results {
			results[i] = &t
		}
		return results, nil
	}

	for i, h := range playerHands {
		if h.Bust() {
			results[i] = &f
			continue
		}

		if h.UpperValue() < dealerHand.UpperValue() {
			results[i] = &f
			continue
		}

		if h.UpperValue() > dealerHand.UpperValue() {
			results[i] = &t
			continue
		}
	}

	return results, nil
}
