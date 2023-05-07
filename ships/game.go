package ships

import (
	battleships "github.com/kovansky/wp-battleships"
	"github.com/rs/zerolog"
)

var _ battleships.Game = (*Game)(nil)

type Game struct {
	key string

	opponent battleships.Player
	player   battleships.Player

	status        battleships.GameStatus
	board         []battleships.Field
	opponentBoard []battleships.Field

	log *zerolog.Logger
}

func NewGame(key string, log *zerolog.Logger) battleships.Game {
	return &Game{key: key, log: log}
}

func (g *Game) SetPlayer(player battleships.Player) {
	g.player = player
}

func (g *Game) Player() battleships.Player {
	return g.player
}

func (g *Game) Key() string {
	return g.key
}

func (g *Game) SetBoard(board []battleships.Field) {
	g.board = board
}

func (g *Game) Board() []battleships.Field {
	return g.board
}

func (g *Game) SetOpponentBoard(board []battleships.Field) {
	g.opponentBoard = board
}

func (g *Game) OpponentBoard() []battleships.Field {
	return g.opponentBoard
}

func (g *Game) SetOpponent(player battleships.Player) {
	g.opponent = player
}

func (g *Game) Opponent() battleships.Player {
	return g.opponent
}

func (g *Game) SetGameStatus(status battleships.GameStatus) {
	g.status = status
}

func (g *Game) GameStatus() battleships.GameStatus {
	return g.status
}
