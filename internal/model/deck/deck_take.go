package deck

import (
	"math/rand"
	"scrub/internal/errors"
	"time"
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

func (d *Deck) TakeRandomCard() (*Card, error) {
	if len(d.ActiveCards) == 0 {
		return nil, errors.ErrActiveCardsIsEmpty
	}

	rand.Seed(time.Now().UnixNano())

	index := rand.Intn(len(d.ActiveCards))

	return d.TakeCardByIndex(index)
}
