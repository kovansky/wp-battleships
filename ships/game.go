package ships

import "github.com/rs/zerolog"

type Game struct {
	Key string
	log *zerolog.Logger
}

func NewGame(key string, log *zerolog.Logger) *Game {
	return &Game{Key: key, log: log}
}
