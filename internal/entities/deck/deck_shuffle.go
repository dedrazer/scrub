package deck

import (
	"math/rand"
	"time"
)

func (d *Deck) Shuffle() {
	d.RestoreBurnt()

	for i := range d.ActiveCards {
		rand.New(rand.NewSource(time.Now().UnixNano()))
		d.Swap(i, rand.Intn(len(d.ActiveCards)))
	}
}
