package battleships

type Game interface {
	Key() string

	SetPlayer(player Player)
	Player() Player
	SetOpponent(player Player)
	Opponent() Player

	SetBoard(board map[string]FieldState)
	Board() map[string]FieldState
	SetOpponentBoard(board map[string]FieldState)
	OpponentBoard() map[string]FieldState

	SetGameStatus(status GameStatus)
	GameStatus() GameStatus
}

type GameStatus struct {
	Status     Status
	LastStatus Status
	ShouldFire bool
	Timer      int
}

type GameUpdateMsg struct{}

type PlayersUpdateMsg struct {
	PlayersInfo string
}

type FieldState string

const (
	FieldStateShip FieldState = "ship"
	FieldStateHit             = "hit"
	FieldStateMiss            = "miss"
	FieldStateSunk            = "sunk"
)

type Field struct {
	Coord string
	State FieldState
}
