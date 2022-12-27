package main

import (
	"blackjack-simulator/internal/service"
	"fmt"
)

func main() {
	testDeck := service.NewDeck()
	card, err := testDeck.GetRandomCard()
	if err != nil {
		fmt.Printf("failed to get card with err: %w", err)
		return
	}

	fmt.Printf("got card: %v", card)
}
