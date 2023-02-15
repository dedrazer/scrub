package deck

import (
	"go.uber.org/zap"
)

func Demo(logger *zap.Logger) {
	logger.Debug("initialising deck")
	testDeck := NewDeck()

	logger.Debug("getting random card")
	card, err := testDeck.GetRandomCard()
	if err != nil {
		logger.Fatal("failed to get card", zap.Error(err))
	}

	card.Log(logger)
	testDeck.Log(logger)

	logger.Debug("taking first card")
	card, err = testDeck.TakeCardByIndex(0)

	card.Log(logger)
	testDeck.Log(logger)

	logger.Debug("taking random card")
	card, err = testDeck.TakeRandomCard()

	card.Log(logger)
	testDeck.Log(logger)

	logger.Debug("taking 10 random cards")
	for i := 0; i < 10; i++ {
		card, err = testDeck.TakeRandomCard()
		card.Log(logger)
	}

	testDeck.Log(logger)

	logger.Debug("getting first card")
	card, _ = testDeck.GetCardByIndex(0)
	card.Log(logger)

	logger.Debug("shuffling cards")
	testDeck.Shuffle()

	logger.Debug("getting first card")
	card, _ = testDeck.GetCardByIndex(0)
	card.Log(logger)
}
