package deck

func (d *Deck) Swap(a, b int) {
	d.ActiveCards[a], d.ActiveCards[b] = d.ActiveCards[b], d.ActiveCards[a]
}
