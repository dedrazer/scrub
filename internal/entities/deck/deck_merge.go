package deck

func (d *Deck) RestoreBurnt() {
	d.ActiveCards = append(d.ActiveCards, d.BurntCards...)
	d.BurntCards = []Card{}
}

func (d *Deck) Merge(other *Deck) {
	d.RestoreBurnt()
	newCards := make([]Card, len(d.ActiveCards)+len(other.ActiveCards)+len(other.BurntCards))
	copy(newCards, d.ActiveCards)
	copy(newCards[len(d.ActiveCards):], other.ActiveCards)
}
