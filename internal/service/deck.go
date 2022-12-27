package service

import "blackjack-simulator/internal/model"

const (
	deckSize = 52
)

var (
	cardSymbols = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	cardValues  = map[string]uint{
		"1":  1,
		"2":  2,
		"3":  3,
		"4":  4,
		"5":  5,
		"6":  6,
		"7":  7,
		"8":  8,
		"9":  9,
		"10": 10,
		"J":  10,
		"Q":  10,
		"K":  10,
		"A":  11,
	}
	suits = []string{"clubs", "diamonds", "hearts", "spades"}
)

func NewDeck() model.Deck {
	res := make([]model.Card, deckSize)
	for i, suit := range suits {
		for j, symbol := range cardSymbols {
			res[i*len(cardSymbols)+j] = NewCard(symbol, suit)
		}
	}

	return model.Deck{ActiveCards: res}
}

func NewCard(symbol, suit string) model.Card {
	return model.Card{
		Value:  cardValues[symbol],
		Symbol: symbol,
		Suit:   suit,
	}
}
