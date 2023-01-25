package player

type Player struct {
	Name    string
	Credits uint64
}

type PlayerBet struct {
	Player    Player
	BetAmount uint64
}
