package model

import (
	"blackjack-simulator/internal/errors"
	"math/rand"
	"time"
)

type Card struct {
	Value  uint
	Symbol string
	Suit   string
}

type Deck struct {
	ActiveCards []Card
	BurntCards  []Card
}

func (d *Deck) GetCardByIndex(index int) (*Card, error) {
	if index >= len(d.ActiveCards) {
		return nil, errors.ErrIndexOutOfRange
	}

	return &d.ActiveCards[index], nil
}

func (d *Deck) GetRandomCard() (*Card, error) {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(d.ActiveCards))

	return d.GetCardByIndex(index)
}
