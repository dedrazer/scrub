package blackjack

import "go.uber.org/zap"

func Demo(logger *zap.Logger) {
	logger.Info("initialising blackjack")
	testBlackjack := NewBlackjack(6)

	var numberOfHands uint8 = 1

	logger.Info("dealing round", zap.Uint8("numberOfHands", numberOfHands))
	playerCards, dealerCards, err := testBlackjack.DealRound(logger, 1)
	if err != nil {
		logger.Fatal("failed to deal round", zap.Error(err))
	}

	for k, v := range playerCards {
		logger.Info("player hand", zap.Uint8("player", k))
		for _, card := range v {
			card.Log(logger)
		}
	}

	logger.Info("dealer hand")
	for _, card := range dealerCards {
		card.Log(logger)
	}
}
