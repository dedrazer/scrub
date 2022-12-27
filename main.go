package main

import (
	"fmt"
	"scrub/internal/model"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to init zap: %w ", err))
	}

	logger.Info("initialising deck")
	testDeck := model.NewDeck()

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
}
