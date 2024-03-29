package deck

import (
	"math/rand"
	"scrub/internal/errorutils"
	"time"
)

func (d *Deck) GetCardByIndex(index int) (*Card, error) {
	if index >= len(d.ActiveCards) {
		return nil, errorutils.ErrIndexOutOfRange
	}

	return &d.ActiveCards[index], nil
}

func (d *Deck) GetRandomCard() (*Card, error) {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(d.ActiveCards))

	return d.GetCardByIndex(index)
}
