package blackjack

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	testBlackjack *Blackjack
	testLogger    *zap.Logger
)

func TestMain(m *testing.M) {
	var err error
	testLogger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	
	resetBJ()

	runCode := m.Run()

	os.Exit(runCode)
}

func resetBJ() {
	testBlackjack = NewBlackjack(testLogger, 10)
}

func TestNewBlackjack(t *testing.T) {
	testLogger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("Failed to start logger: %s", err.Error())
	}

	for i := 1; i <= 10; i++ {
		b := NewBlackjack(testLogger, uint(i))

		assert.Len(t, b.deck.ActiveCards, 52*i, "%d deck(s) must have %d cards", i, i*52)
	}
}
