package blackjack

import "go.uber.org/zap"

type DealerHand struct {
	Hand
}

func (dh *DealerHand) DealerLog(logger *zap.Logger) {
	logger.Debug("dealer hand", zap.String("card", dh.cards[1].Print()))
}

func (dh *DealerHand) DealerResult(logger *zap.Logger) {
	dh.Log(logger, "dealer result")
}
func (dh *DealerHand) UpCardValue() uint {
	return dh.cards[1].Value
}
