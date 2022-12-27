package model

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	descriptionCard = "%s of %s"
	descriptionDeck = "Deck with %d remaining cards and %d burnt cards"
)

type Card struct {
	Value  uint
	Symbol string
	Suit   string
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

func (d *Deck) Print() string {
	return fmt.Sprintf(descriptionDeck, len(d.ActiveCards), len(d.BurntCards))
}

func (d *Deck) Log(logger *zap.Logger) {
	logger.Info("deck info", zap.String("deck", d.Print()))
}
