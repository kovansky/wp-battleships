package battleships

type Game interface {
	Key() string

	SetPlayer(player Player)
	Player() Player
	SetOpponent(player Player)
	Opponent() Player

	SetBoard(board []string)
	Board() []string

	SetGameStatus(status GameStatus)
	GameStatus() GameStatus
}

type GameStatus struct {
	Status     string
	LastStatus string
	ShouldFire bool
	Timer      int
}
