package deck

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	deckSize        = 52
	descriptionCard = "%s of %s"
	descriptionDeck = "Deck with %d remaining cards and %d burnt cards"
)

var (
	cardSymbols = []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}
	CardValues  = map[string]uint{
		"Ace":   1,
		"2":     2,
		"3":     3,
		"4":     4,
		"5":     5,
		"6":     6,
		"7":     7,
		"8":     8,
		"9":     9,
		"10":    10,
		"Jack":  10,
		"Queen": 10,
		"King":  10,
	}
	suits = []string{"Clubs", "Diamonds", "Hearts", "Spades"}
)

type Card struct {
	Value  uint
	Symbol string
	Suit   string
}

func NewCard(symbol, suit string) Card {
	return Card{
		Value:  CardValues[symbol],
		Symbol: symbol,
		Suit:   suit,
	}
}

func (c *Card) Print() string {
	return fmt.Sprintf(descriptionCard, c.Symbol, c.Suit)
}

func (c *Card) Log(logger *zap.Logger) {
	logger.Info("got card", zap.String("card", c.Print()))
}

type Deck struct {
	ActiveCards []Card
	BurntCards  []Card
}

func NewDeck() Deck {
	res := make([]Card, deckSize)
	for i, suit := range suits {
		for j, symbol := range cardSymbols {
			res[i*len(cardSymbols)+j] = NewCard(symbol, suit)
		}
	}

	return Deck{ActiveCards: res, BurntCards: []Card{}}
}

func NewDeckByCards(cards []Card) Deck {
	return Deck{ActiveCards: cards, BurntCards: []Card{}}
}

func (d *Deck) Print() string {
	return fmt.Sprintf(descriptionDeck, len(d.ActiveCards), len(d.BurntCards))
}

func (d *Deck) Log(logger *zap.Logger) {
	logger.Info("deck info", zap.String("deck", d.Print()))
}
