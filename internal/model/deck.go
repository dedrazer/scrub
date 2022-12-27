package model

import "blackjack-simulator/internal/errors"

type Card struct {
	Value  uint
	Symbol string
	Suit   string
}

type Deck struct {
	ActiveCards []Card
	BurntCards  []Card
}

func (d *Deck) GetCard() (*Card, error) {
	if len(d.ActiveCards) > 0 {
		res := d.ActiveCards[0]
		d.ActiveCards = d.ActiveCards[1:]
		d.BurntCards = append(d.BurntCards, res)
		return &res, nil
	}

	return nil, errors.ErrActiveCardsIsEmpty
}
