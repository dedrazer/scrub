package player

import "scrub/internal/errors"

type Player struct {
	Name    string
	Credits uint64
}

func (p *Player) Win(amount uint64) {
	p.Credits += amount
}

func (p *Player) Lose(amount uint64) error {
	if p.Credits < amount {
		return errors.ErrInsufficientCredits
	}

	p.Credits -= amount
	return nil
}
