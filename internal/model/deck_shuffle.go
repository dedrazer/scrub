package model

import (
	"math/rand"
	"time"
)

func (d *Deck) Shuffle() {
	d.RestoreBurnt()

	for i := range d.ActiveCards {
		rand.Seed(time.Now().UnixNano())
		d.Swap(i, rand.Intn(len(d.ActiveCards)))
	}
}
