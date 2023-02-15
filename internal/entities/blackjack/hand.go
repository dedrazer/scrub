package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	"strings"

	"go.uber.org/zap"
)

type Hand struct {
	cards     []deck.Card
	isSplit   bool
	isDoubled bool
	result    *string
	BetAmount uint64
}

func (h *Hand) Print() string {
	cardStrings := make([]string, len(h.cards))
	for i, c := range h.cards {
		cardStrings[i] = c.Print()
	}

	return strings.Join(cardStrings, ", ")
}

func (h *Hand) Log(logger *zap.Logger, args ...string) {
	if len(h.cards) == 0 {
		logger.Debug("hand", zap.String("cards", "empty"))
		return
	}

	if len(args) == 0 {
		logger.Debug("hand", zap.String("cards", h.Print()), zap.String("value", h.PrintValue()))
		return
	}

	logger.Debug(args[0], zap.String("cards", h.Print()), zap.String("value", h.PrintValue()))
}

func (h *Hand) PrintValue() string {
	values := h.Value()

	stringValues := make([]string, len(values))
	for i, v := range values {
		stringValues[i] = fmt.Sprintf("%d", v)
	}

	return strings.Join(stringValues, "/")
}

func (h *Hand) LogValue(logger *zap.Logger) {
	logger.Debug("hand value", zap.String("amount", h.PrintValue()))
}

func (h *Hand) Value() []uint {
	var total uint

	containsAce := false
	for _, c := range h.cards {
		if c.Value == 1 {
			if containsAce && total == 1 {
				total += 10
			}

			containsAce = true
		}

		total += c.Value
	}

	if total < 12 && containsAce {
		return []uint{total, total + 10}
	}

	return []uint{total}
}

func (h *Hand) UpperValue() uint {
	values := h.Value()
	return values[len(values)-1]
}

func (h *Hand) Bust() bool {
	return h.UpperValue() > 21
}

func (h *Hand) AddCard(c deck.Card) {
	if h.cards == nil {
		h.cards = []deck.Card{}
	}
	h.cards = append(h.cards, c)
}

func (h *Hand) Blackjack() bool {
	return len(h.cards) == 2 && h.UpperValue() == 21 && !h.isSplit
}

func (h *Hand) CanSplit() bool {
	return !h.isSplit && len(h.cards) == 2 && h.cards[0].Symbol == h.cards[1].Symbol
}

func (h *Hand) CanDouble() bool {
	return len(h.cards) == 2 && !h.isSplit
}

func (h *Hand) IsSoft() bool {
	return len(h.Value()) > 1
}

func (h *Hand) Double() {
	h.BetAmount *= 2
	h.isDoubled = true
}

func (h *Hand) ResetDouble() {
	if h.isDoubled {
		h.BetAmount /= 2
		h.isDoubled = false
	}
}
