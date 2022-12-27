package model

import (
	"blackjack-simulator/internal/errors"
)

func (d *Deck) TakeCardByIndex(index int) (*Card, error) {
	if index >= len(d.ActiveCards) {
		return nil, errors.ErrIndexOutOfRange
	}

	res := d.ActiveCards[index]
	d.ActiveCards = append(d.ActiveCards[:index], d.ActiveCards[index+1:]...)
	d.BurntCards = append(d.BurntCards, res)

	return &d.BurntCards[len(d.BurntCards)-1], nil
}
