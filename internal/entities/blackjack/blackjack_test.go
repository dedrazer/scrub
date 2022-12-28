package blackjack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlackjack(t *testing.T) {
	for i := 1; i < 10; i++ {
		b := NewBlackjack(uint(i))

		assert.Len(t, b.deck.ActiveCards, 52*i, "%d deck(s) must have %d cards", i, i*52)
	}
}
