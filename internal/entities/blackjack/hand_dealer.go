package blackjack

import "go.uber.org/zap"

type DealerHand struct {
	Hand
}

func (dh *DealerHand) shouldDraw() bool {
	return dh.UpperValue() < 17
}

func (dh *DealerHand) hasSoftValue() bool {
	return len(dh.Value()) > 1
}

func (dh *DealerHand) hasNoValue() bool {
	return len(dh.Value()) == 0 || dh.UpperValue() == 0
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
