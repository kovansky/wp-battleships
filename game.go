package battleships

type Game interface {
	Key() string

	SetPlayer(player Player)
	Player() Player
	SetOpponent(player Player)
	Opponent() Player

	SetBoard(board []Field)
	Board() []Field
	SetOpponentBoard(board []Field)
	OpponentBoard() []Field

	SetGameStatus(status GameStatus)
	GameStatus() GameStatus
}

type GameStatus struct {
	Status     Status
	LastStatus Status
	ShouldFire bool
	Timer      int
}

type GameUpdateMsg struct {
	Board         []Field
	OpponentBoard []Field
}

type BoardUpdateMsg struct {
	Board []Field
}

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
