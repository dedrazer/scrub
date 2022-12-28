package deck

import (
	"go.uber.org/zap"
)

func Demo(logger *zap.Logger) {
	logger.Info("initialising deck")
	testDeck := NewDeck()

	logger.Info("getting random card")
	card, err := testDeck.GetRandomCard()
	if err != nil {
		logger.Fatal("failed to get card", zap.Error(err))
	}

	card.Log(logger)
	testDeck.Log(logger)

	logger.Info("taking first card")
	card, err = testDeck.TakeCardByIndex(0)

	card.Log(logger)
	testDeck.Log(logger)

	logger.Info("taking random card")
	card, err = testDeck.TakeRandomCard()

	card.Log(logger)
	testDeck.Log(logger)

	logger.Info("taking 10 random cards")
	for i := 0; i < 10; i++ {
		card, err = testDeck.TakeRandomCard()
		card.Log(logger)
	}

	testDeck.Log(logger)

	logger.Info("getting first card")
	card, _ = testDeck.GetCardByIndex(0)
	card.Log(logger)

	logger.Info("shuffling cards")
	testDeck.Shuffle()

	logger.Info("getting first card")
	card, _ = testDeck.GetCardByIndex(0)
	card.Log(logger)
}
