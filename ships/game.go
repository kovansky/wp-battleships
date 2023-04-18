package ships

import (
	battleships "github.com/kovansky/wp-battleships"
	"github.com/rs/zerolog"
)

var _ battleships.Game = (*Game)(nil)

type Game struct {
	key    string
	player battleships.Player

	log *zerolog.Logger
}

func NewGame(key string, log *zerolog.Logger) battleships.Game {
	return &Game{key: key, log: log}
}

func (g Game) SetPlayer(player battleships.Player) {
	g.player = player
}

func (g Game) Player() battleships.Player {
	return g.player
}

func (g Game) Key() string {
	return g.key
}
