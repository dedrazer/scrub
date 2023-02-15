package player

import (
	"fmt"
	"scrub/internal/errors"
)

type Player struct {
	Name    string
	Credits uint64
	Wins    uint64
	Losses  uint64
}

func (p *Player) Win(amount uint64) {
	p.Credits += amount

	p.Wins++
}

func (p *Player) Lose(amount uint64) error {
	if p.Credits < amount {
		return errors.ErrInsufficientCredits
	}

	p.Credits -= amount

	p.Losses++
	return nil
}

func (p *Player) WinRate() float64 {
	return float64(p.Wins) / float64(p.Wins+p.Losses)
}

func (p *Player) WinRateString() string {
	return fmt.Sprintf("%.2f%%", p.WinRate()*100)
}
