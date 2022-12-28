package main

import (
	"fmt"
	"scrub/internal/entities/deck"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to init zap: %w ", err))
	}

	logger.Info("initialising deck")
	testDeck := deck.NewDeck()

	logger.Info("getting random card")
	card, err := testDeck.GetRandomCard()
	if err != nil {
		fmt.Printf("failed to get card with err: %w", err)
		return
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
