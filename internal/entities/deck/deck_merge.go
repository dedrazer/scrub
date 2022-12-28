package deck

func (d *Deck) RestoreBurnt() {
	d.ActiveCards = append(d.ActiveCards, d.BurntCards...)
	d.BurntCards = []Card{}
}

func Merge(deck1, deck2 *Deck) Deck {
	deck1.RestoreBurnt()
	deck2.RestoreBurnt()

	newCards := make([]Card, len(deck1.ActiveCards)+len(deck2.ActiveCards))
	copy(newCards, deck1.ActiveCards)
	copy(newCards[len(deck1.ActiveCards):], deck2.ActiveCards)

	return NewDeckByCards(newCards)
}
